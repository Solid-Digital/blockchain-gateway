package pipeline

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	stderr "errors"
	"fmt"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	v1 "k8s.io/api/apps/v1"

	"github.com/BurntSushi/toml"
	"github.com/volatiletech/null"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/pkg/pipeline/internal/janus/pkg/janus"
	"bitbucket.org/unchain/ares/pkg/pipeline/internal/janus/pkg/pipeline"
	"github.com/unchainio/pkg/xlogger"

	"github.com/unchainio/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"

	"bitbucket.org/unchain/ares/gen/dto"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) DeployConfiguration(params *dto.DeployConfigurationRequest, orgName string, pipelineName string, envName string, user *dto.User) (*dto.GetDeploymentResponse, *apperr.Error) {
	var org *orm.Organization
	var pipeline *orm.Pipeline
	var env *orm.Environment
	var deployment *orm.Deployment
	var kubeDeployment *v1.Deployment

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var err error
		var appErr *apperr.Error

		s.log.Debugf("Deploying configuration")
		org, pipeline, appErr = xorm.GetPpelineTx(ctx, tx, orgName, pipelineName)
		if appErr != nil {
			return appErr
		}

		var config *orm.Configuration
		mods := []qm.QueryMod{
			qm.Load(orm.ConfigurationRels.CreatedBy),
			qm.Load(orm.ConfigurationRels.UpdatedBy),
			qm.Load(qm.Rels(
				orm.ConfigurationRels.BaseConfiguration,
				orm.BaseConfigurationRels.Version,
				orm.BaseVersionRels.Base,
			)),
			qm.Load(qm.Rels(
				orm.ConfigurationRels.TriggerConfiguration,
				orm.TriggerConfigurationRels.Version,
				orm.TriggerVersionRels.Trigger,
			)),
			qm.Load(qm.Rels(
				orm.ConfigurationRels.ActionConfigurations,
			),
				qm.OrderBy(orm.ActionConfigurationColumns.Index),
			),
			qm.Load(qm.Rels(
				orm.ConfigurationRels.ActionConfigurations,
				orm.ActionConfigurationRels.Version,
				orm.ActionVersionRels.Action,
			)),
		}
		// if the desired revision is -1, get the latest revision
		if *params.ConfigurationRevision == -1 {
			mods = append(mods, qm.OrderBy(orm.ConfigurationColumns.Revision+" DESC"))
		} else {
			mods = append(mods, orm.ConfigurationWhere.Revision.EQ(*params.ConfigurationRevision))
		}

		config, err = pipeline.Configurations(mods...).One(ctx, tx)
		if err != nil {
			err := ares.ParsePQErr(err)
			switch {
			case stderr.Is(err, apperr.NotFound):
				return ares.ErrPipelineConfigNotFound(err, orgName, pipelineName, *params.ConfigurationRevision)
			default:
				return err
			}
		}

		params.ConfigurationRevision = &config.Revision

		env, appErr := xorm.GetOrgEnvironmentTx(ctx, tx, org, envName)
		if appErr != nil {
			return appErr
		}

		replicas, appErr := getReplicaCount(params, env)
		if appErr != nil {
			return appErr
		}

		deployment = &orm.Deployment{
			PipelineID:      pipeline.ID,
			ConfigurationID: config.ID,
			CreatedByID:     user.ID,
			UpdatedByID:     user.ID,
			EnvironmentID:   env.ID,
			Replicas:        replicas,
			FullName:        getDeploymentFullName(pipeline.Name, envName),
			URL:             fmt.Sprintf("https://%s.%s.%s/%s--%s", orgName, s.cfg.DeploymentRegion, s.cfg.DeploymentHost, pipeline.Name, env.Name),
			Image:           fmt.Sprintf("%s/%s/%s:%d", s.cfg.RegistryURL, orgName, pipeline.Name, config.Revision),
			Host:            fmt.Sprintf("%s.%s.%s", orgName, s.cfg.DeploymentRegion, s.cfg.DeploymentHost),
			Path:            fmt.Sprintf("/%s--%s(/|$)(.*)", pipeline.Name, env.Name),
			RewriteTarget:   "/$2",
			Dirty:           false,
		}
		s.log.Debugf("Getting config file")

		configString, appErr := s.getConfig(config)
		if appErr != nil {
			return appErr
		}

		s.log.Debugf("Checking if organization should be created in Harbor")
		// If harbor was configured, upsert a project for this organization.
		if s.service.registry != nil {
			s.log.Debugf("Registry is set, inserting organization")
			_, err = s.service.registry.UpsertProject(orgName, false)
			if err != nil {
				return apperr.Internal.Wrap(errors.Wrap(err, ""))
			}
		}

		s.log.Debugf("Getting all components")
		components := s.getAllComponents(config)

		s.log.Debugf("Building image")
		err = s.service.imageBuilder.BuildImage(&ares.BuildManifest{
			Tag:        deployment.Image,
			BaseImage:  config.R.BaseConfiguration.R.Version.DockerImageRef,
			Components: components,
			Cmd:        []string{config.R.BaseConfiguration.R.Version.Entrypoint, "--cfg=/etc/opt/unchain/config.toml"},
		})
		if err != nil {
			return apperr.Internal.Wrap(errors.Wrap(err, ""))
		}

		err = deployment.Upsert(ctx, tx, true, []string{orm.DeploymentColumns.EnvironmentID, orm.DeploymentColumns.PipelineID}, boil.Infer(), boil.Infer())
		if err != nil {
			return ares.ParsePQErr(err)
		}

		s.log.Debugf("Deploying to kubernetes")

		envVars, appErr := loadEnvVarsTx(ctx, tx, orgName, pipelineName, envName)
		if appErr != nil {
			return appErr
		}

		envVarDTOs := make([]*orm.EnvironmentVariableDTO, len(envVars))

		for i, v := range envVars {
			envVarDTOs[i] = v.DTO()
		}

		params := &ares.DeploymentParams{
			DeploymentDTO:        deployment.DTO(),
			Config:               configString,
			OrganizationName:     orgName,
			EnvironmentName:      envName,
			Revision:             *params.ConfigurationRevision,
			EnvironmentVariables: envVarDTOs,
		}

		kubeDeployment, err = s.service.kube.CreateDeployment(params)
		if err != nil {
			return apperr.Internal.Wrap(errors.Wrap(err, ""))
		}

		// mark env vars as deployed
		for _, v := range envVars {
			v.Deployed = true
			v.UpdatedByID = user.ID

			_, err = v.Update(ctx, tx, boil.Infer())
			if err != nil {
				return ares.ParsePQErr(err)
			}
		}

		err = deployment.L.LoadConfiguration(ctx, tx, true, deployment, nil)
		if err != nil {
			return ares.ParsePQErr(err)
		}

		err = deployment.L.LoadCreatedBy(ctx, tx, true, deployment, nil)
		if err != nil {
			return ares.ErrLoadingCreatedBy(ares.ParsePQErr(err), "deployment")
		}

		err = deployment.L.LoadUpdatedBy(ctx, tx, true, deployment, nil)
		if err != nil {
			return ares.ErrLoadingUpdatedBy(ares.ParsePQErr(err), "deployment")
		}

		return nil
	})

	if appErr != nil {
		return nil, appErr
	}

	return getDeployment(pipeline, env, deployment, kubeDeployment), nil
}

func getDeploymentFullName(pipelineName string, envName string) string {
	return fmt.Sprintf("%s-%s", pipelineName, envName)
}

func (s *Service) getConfig(config *orm.Configuration) (string, *apperr.Error) {
	var appErr *apperr.Error

	cfg, appErr := getAdapterConfig(config)
	if appErr != nil {
		return "", appErr
	}

	// it doesn't accept pointers, don't remove the *
	buf := bytes.NewBuffer(nil)
	err := toml.NewEncoder(buf).Encode(*cfg)

	if err != nil {
		return "", apperr.Internal.Wrap(errors.Wrap(err, ""))
	}

	return buf.String(), nil
}

func getAdapterConfig(config *orm.Configuration) (*janus.Config, *apperr.Error) {
	triggerMessageConfig, err := getAdapterMessageConfigObject(config.R.TriggerConfiguration.MessageConfig) // TODO(e-nikolov) test me
	if err != nil {
		return nil, apperr.Internal.Wrap(err)
	}

	trigger := pipeline.ComponentConfig{
		Name:          config.R.TriggerConfiguration.Name,
		Path:          config.R.TriggerConfiguration.R.Version.FileName,
		Config:        config.R.TriggerConfiguration.Config,
		MessageConfig: triggerMessageConfig,
	}

	actions := make([]pipeline.ComponentConfig, len(config.R.ActionConfigurations))

	for i, action := range config.R.ActionConfigurations {
		actionMessageConfig, err := getAdapterMessageConfigObject(action.MessageConfig)
		if err != nil {
			return nil, apperr.Internal.Wrap(err)
		}

		actions[i] = pipeline.ComponentConfig{
			Name:          action.Name,
			Path:          action.R.Version.FileName,
			Config:        action.Config,
			MessageConfig: actionMessageConfig,
		}
	}

	return &janus.Config{
		Logger: &xlogger.Config{
			Level:  "debug",
			Format: "json",
		},
		Pipeline: &pipeline.Config{
			Trigger: trigger,
			Actions: actions,
		},
	}, nil
}

func getAdapterMessageConfigObject(cfg null.JSON) (map[string]interface{}, error) {
	res := make(map[string]interface{})

	if cfg.Valid {
		err := json.Unmarshal(cfg.JSON, &res)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
	}

	return res, nil
}

func getAdapterMessageConfigTOML(cfg null.JSON) (string, *apperr.Error) {
	if cfg.Valid {
		tomlBytes, err := jsonToTOML(cfg.JSON)
		if err != nil {
			return "", apperr.Internal.Wrap(err)
		}

		return string(tomlBytes), nil
	}

	return "", nil
}

func getReplicaCount(params *dto.DeployConfigurationRequest, env *orm.Environment) (int64, *apperr.Error) {
	if params.Replicas == nil {
		return env.MaxReplicas, nil
	}

	if *params.Replicas > env.MaxReplicas {
		return 0, ares.ErrTooManyReplicas(apperr.BadRequest.Wrap(errors.New("")), *params.Replicas, env.MaxReplicas)
	}

	return *params.Replicas, nil
}

func (s *Service) getAllComponents(config *orm.Configuration) []*ares.Component {
	capacity := len(config.R.ActionConfigurations) + 1 // 1 = trigger
	components := make([]*ares.Component, capacity)

	components[0] = &ares.Component{
		FileName: config.R.TriggerConfiguration.R.Version.FileName,
		FileId:   config.R.TriggerConfiguration.R.Version.FileID,
	}

	for i, action := range config.R.ActionConfigurations {
		components[i+1] = &ares.Component{
			FileName: action.R.Version.FileName,
			FileId:   action.R.Version.FileID,
		}
	}

	return components
}
