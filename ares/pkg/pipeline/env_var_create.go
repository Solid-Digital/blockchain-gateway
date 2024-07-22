package pipeline

import (
	"context"
	"database/sql"
	stderr "errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/volatiletech/sqlboiler/boil"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) CreateEnvironmentVariable(params *dto.CreateEnvironmentVariableRequest, orgName string, pipelineName string, envName string, user *dto.User) (*dto.GetEnvironmentVariableResponse, *apperr.Error) {
	var ret *dto.GetEnvironmentVariableResponse
	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error
		var err error

		org, pipeline, appErr := xorm.GetPpelineTx(ctx, tx, orgName, pipelineName)
		if appErr != nil {
			return appErr
		}

		env, appErr := xorm.GetOrgEnvironmentTx(ctx, tx, org, envName)
		if appErr != nil {
			return appErr
		}

		var index int64
		topEnvVar, err := pipeline.EnvironmentVariables(orm.EnvironmentVariableWhere.EnvironmentID.EQ(env.ID), qm.OrderBy(orm.EnvironmentVariableColumns.Index+" DESC")).One(ctx, tx)
		appErr = ares.ParsePQErr(err)
		switch {
		case appErr == nil:
			index = topEnvVar.Index + 1
		case stderr.Is(appErr, apperr.NotFound):
			index = 1
		default:
			return appErr
		}

		envVar := &orm.EnvironmentVariable{
			OrganizationID: org.ID,
			PipelineID:     pipeline.ID,
			EnvironmentID:  env.ID,
			CreatedByID:    user.ID,
			UpdatedByID:    user.ID,
			Index:          index,
			Key:            params.Key,
			Value:          params.Value,
			Secret:         params.Secret,
			Deployed:       false,
		}

		err = envVar.Insert(ctx, tx, boil.Infer())
		if err != nil {
			err := ares.ParsePQErr(err)
			switch {
			case stderr.Is(err, apperr.Conflict):
				return ares.ErrDuplicateEnvVar(err, orgName, pipelineName, envName, params.Key)
			default:
				return err
			}
		}

		appErr = markExistingDeploymentDirtyTx(ctx, tx, true, org, pipeline, env)
		if appErr != nil {
			return appErr
		}

		err = envVar.L.LoadCreatedBy(ctx, tx, true, envVar, nil)
		if err != nil {
			return ares.ErrLoadingCreatedBy(ares.ParsePQErr(err), params.Key)
		}

		err = envVar.L.LoadUpdatedBy(ctx, tx, true, envVar, nil)
		if err != nil {
			return ares.ErrLoadingUpdatedBy(ares.ParsePQErr(err), params.Key)
		}

		ret = getEnvironmentVariable(envVar)

		return nil
	})

	if appErr != nil {
		return nil, appErr
	}

	return ret, nil
}

func markExistingDeploymentDirtyTx(ctx context.Context, tx *sql.Tx, dirty bool, org *orm.Organization, pipeline *orm.Pipeline, env *orm.Environment) *apperr.Error {
	deployment, appErr := xorm.GetDeploymentTx(ctx, tx, org, pipeline, env,
		qm.Load(orm.DeploymentRels.Configuration),
		qm.Load(orm.DeploymentRels.CreatedBy),
		qm.Load(orm.DeploymentRels.UpdatedBy),
	)

	if appErr != nil {
		switch {
		case stderr.Is(appErr, apperr.NotFound):
			return nil
		default:
			return appErr
		}
	}

	deployment.Dirty = dirty
	_, err := deployment.Update(ctx, tx, boil.Infer())
	if err != nil {
		return ares.ParsePQErr(err)
	}

	return nil
}
