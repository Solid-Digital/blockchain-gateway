package pipeline

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	v1 "k8s.io/api/apps/v1"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/go-openapi/strfmt"

	"bitbucket.org/unchain/ares/gen/dto"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) GetPipeline(orgName string, pipelineName string) (*dto.GetPipelineResponse, *apperr.Error) {
	var ret *dto.GetPipelineResponse

	err := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		org, pipeline, appErr := xorm.GetPpelineTx(ctx, tx, orgName, pipelineName,
			qm.Load(orm.PipelineRels.CreatedBy),
			qm.Load(orm.PipelineRels.UpdatedBy),
			qm.Load(orm.PipelineRels.Deployments),
			qm.Load(qm.Rels(orm.PipelineRels.EnvironmentVariables, orm.EnvironmentVariableRels.CreatedBy)),
			qm.Load(qm.Rels(orm.PipelineRels.EnvironmentVariables, orm.EnvironmentVariableRels.UpdatedBy)),
			qm.Load(qm.Rels(orm.PipelineRels.Deployments, orm.DeploymentRels.Environment)),
			qm.Load(orm.PipelineRels.Configurations),
			qm.Load(qm.Rels(
				orm.PipelineRels.Deployments,
				orm.DeploymentRels.CreatedBy,
			)),
			qm.Load(qm.Rels(
				orm.PipelineRels.Deployments,
				orm.DeploymentRels.UpdatedBy,
			)),
			qm.Load(qm.Rels(
				orm.PipelineRels.Deployments,
				orm.DeploymentRels.Configuration,
				orm.ConfigurationRels.BaseConfiguration,
				orm.BaseConfigurationRels.Version,
				orm.BaseVersionRels.Base,
			)),
			qm.Load(qm.Rels(
				orm.PipelineRels.Deployments,
				orm.DeploymentRels.Configuration,
				orm.ConfigurationRels.TriggerConfiguration,
				orm.TriggerConfigurationRels.Version,
				orm.TriggerVersionRels.Trigger,
			)),
			qm.Load(qm.Rels(
				orm.PipelineRels.Deployments,
				orm.DeploymentRels.Configuration,
				orm.ConfigurationRels.CreatedBy,
			)),
			qm.Load(qm.Rels(
				orm.PipelineRels.Deployments,
				orm.DeploymentRels.Configuration,
				orm.ConfigurationRels.UpdatedBy,
			)),
			qm.Load(qm.Rels(
				orm.PipelineRels.Deployments,
				orm.DeploymentRels.Configuration,
				orm.ConfigurationRels.ActionConfigurations,
			),
				qm.OrderBy(orm.ActionConfigurationColumns.Index),
			),
			qm.Load(qm.Rels(
				orm.PipelineRels.Deployments,
				orm.DeploymentRels.Configuration,
				orm.ConfigurationRels.ActionConfigurations,
				orm.ActionConfigurationRels.Version,
				orm.ActionVersionRels.Action,
			)),
		)
		if appErr != nil {
			return appErr
		}

		envs, appErr := xorm.GetAllOrgEnvironmentsTx(ctx, tx, org)
		if appErr != nil {
			return appErr
		}

		ret, appErr = s.toDTOPipeline(orgName, pipeline, envs)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (s *Service) toDTOPipeline(orgName string, p *orm.Pipeline, envs orm.EnvironmentSlice) (*dto.GetPipelineResponse, *apperr.Error) {
	environments, err := s.getPipelineTxEnvironments(orgName, p, envs, p.R.Deployments)

	if err != nil {
		return nil, err
	}

	return &dto.GetPipelineResponse{
		ID:                     &p.ID,
		DisplayName:            &p.DisplayName,
		Name:                   &p.Name,
		Description:            &p.Description,
		Status:                 &p.Status,
		Environments:           environments,
		ConfigurationRevisions: getPipelineTxConfigurations(p.R.Configurations),
		CreatedAt:              (*strfmt.DateTime)(&p.CreatedAt),
		CreatedBy: &dto.CreatedBy{
			ID:       p.R.CreatedBy.ID,
			FullName: p.R.CreatedBy.FullName,
		},
		UpdatedAt: (*strfmt.DateTime)(&p.UpdatedAt),
		UpdatedBy: &dto.UpdatedBy{
			FullName: p.R.UpdatedBy.FullName,
			ID:       p.R.UpdatedBy.ID,
		},
	}, nil
}

func getPipelineTxConfigurations(configs orm.ConfigurationSlice) []*dto.PipelineConfigurationRevision {
	res := make([]*dto.PipelineConfigurationRevision, len(configs))

	for i, config := range configs {
		res[i] = &dto.PipelineConfigurationRevision{
			Revision:      config.Revision,
			CommitMessage: config.CommitMessage,
		}
	}

	return res
}

func (s *Service) getPipelineTxEnvironments(orgName string, pipeline *orm.Pipeline, envs orm.EnvironmentSlice, deployments orm.DeploymentSlice) ([]*dto.PipelineEnvironment, *apperr.Error) {
	dd := make([]*dto.PipelineEnvironment, len(envs))

	envMap := make(map[int64]int)

	for i, env := range envs {
		envMap[env.ID] = i
		envVars := s.getPipelineTxEnvironmentsVars(pipeline, env)

		dd[i] = &dto.PipelineEnvironment{
			Name:      env.Name,
			Variables: envVars,
		}
	}

	kubeDeployments, appErr := s.getKubeDeployments(orgName, deployments)
	if appErr != nil {
		return nil, appErr
	}

	for _, d := range deployments {
		kd, ok := kubeDeployments[d.FullName]
		var availableReplicas int64 = 0

		if ok {
			availableReplicas = int64(kd.Status.AvailableReplicas)
		}

		dd[envMap[d.EnvironmentID]] = &dto.PipelineEnvironment{
			Deployment: &dto.GetDeploymentResponse{
				ID:        &d.ID,
				CreatedAt: (*strfmt.DateTime)(&d.CreatedAt),
				CreatedBy: &dto.CreatedBy{
					ID:       d.R.CreatedBy.ID,
					FullName: d.R.CreatedBy.FullName,
				},
				UpdatedAt: (*strfmt.DateTime)(&d.UpdatedAt),
				UpdatedBy: &dto.UpdatedBy{
					ID:       d.R.UpdatedBy.ID,
					FullName: d.R.UpdatedBy.FullName,
				},
				DesiredReplicas:   &d.Replicas,
				AvailableReplicas: &availableReplicas,
				Configuration: &dto.PipelineConfigurationRevision{
					CommitMessage: d.R.Configuration.CommitMessage,
					Revision:      d.R.Configuration.Revision,
				},
				Dirty: nil,
				Host:  &d.Host,
				Image: &d.Image,
				Path:  &d.Path,
				URL:   &d.URL,
			},
			Name:      d.R.Environment.Name,
			Variables: dd[envMap[d.EnvironmentID]].Variables,
		}
	}

	return dd, nil
}

func (s *Service) getPipelineTxEnvironmentsVars(pipeline *orm.Pipeline, env *orm.Environment) dto.GetAllEnvironmentVariablesResponse {
	envVarsRes := make([]*dto.GetEnvironmentVariableResponse, 0)
	for _, v := range pipeline.R.EnvironmentVariables {
		if v.EnvironmentID == env.ID {
			envVarsRes = append(envVarsRes, getEnvironmentVariable(v))
		}
	}

	return envVarsRes
}

func (s *Service) getKubeDeployments(orgName string, deployments orm.DeploymentSlice) (map[string]*v1.Deployment, *apperr.Error) {
	m := make(map[string]*v1.Deployment)
	if len(deployments) == 0 {
		return m, nil
	}

	var names []string
	for _, deployment := range deployments {
		names = append(names, deployment.FullName)
	}

	kubeDeployments, err := s.service.kube.GetAllDeployments(orgName, names)
	if err != nil {
		return nil, apperr.Internal.Wrap(err).WithMessage("failed to fetch all pipeline deployments")
	}

	for _, kubeDeployment := range kubeDeployments.Items {
		m[kubeDeployment.Name] = &kubeDeployment
	}

	return m, nil
}
