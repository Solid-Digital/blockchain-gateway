package pipeline

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/volatiletech/sqlboiler/boil"
)

func (s *Service) UpdatePipeline(params *dto.UpdatePipelineRequest, orgName string, pipelineName string, principal *dto.User) (*dto.GetPipelineResponse, *apperr.Error) {
	var ret *dto.GetPipelineResponse

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

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

		pipeline.UpdatedByID = principal.ID
		pipeline.DisplayName = params.DisplayName
		pipeline.Description = params.Description

		_, err := pipeline.Update(ctx, tx,
			boil.Whitelist(
				orm.PipelineColumns.UpdatedByID,
				orm.PipelineColumns.DisplayName,
				orm.PipelineColumns.Description,
			))
		if err != nil {
			return ares.ParsePQErr(err)
		}

		err = pipeline.L.LoadCreatedBy(ctx, tx, true, pipeline, nil)
		if err != nil {
			return ares.ErrLoadingCreatedBy(ares.ParsePQErr(err), pipelineName)
		}

		err = pipeline.L.LoadUpdatedBy(ctx, tx, true, pipeline, nil)
		if err != nil {
			return ares.ErrLoadingUpdatedBy(ares.ParsePQErr(err), pipelineName)
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
	if appErr != nil {
		return nil, appErr
	}

	return ret, nil
}
