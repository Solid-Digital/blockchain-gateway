package pipeline

import (
	"context"
	"database/sql"
	stderr "errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/go-openapi/strfmt"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) GetAllEnvironmentVariables(orgName string, pipelineName string, envName string, user *dto.User) (dto.GetAllEnvironmentVariablesResponse, *apperr.Error) {
	var ret dto.GetAllEnvironmentVariablesResponse

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		envVars, appErr := loadEnvVarsTx(ctx, tx, orgName, pipelineName, envName)
		if appErr != nil {
			return appErr
		}

		ret = make(dto.GetAllEnvironmentVariablesResponse, len(envVars))
		for i, envVar := range envVars {
			ret[i] = getEnvironmentVariable(envVar)
		}

		return nil
	})

	if appErr != nil {
		return nil, appErr
	}

	return ret, nil
}

func loadEnvVarsTx(ctx context.Context, tx *sql.Tx, orgName string, pipelineName string, envName string) (orm.EnvironmentVariableSlice, *apperr.Error) {
	var appErr *apperr.Error
	org, pipeline, appErr := xorm.GetPpelineTx(ctx, tx, orgName, pipelineName)
	if appErr != nil {
		return nil, appErr
	}

	env, appErr := xorm.GetOrgEnvironmentTx(ctx, tx, org, envName)
	if appErr != nil {
		return nil, appErr
	}

	envVars, err := pipeline.EnvironmentVariables(
		orm.EnvironmentVariableWhere.EnvironmentID.EQ(env.ID),
		qm.OrderBy(orm.EnvironmentVariableColumns.Index+" DESC"),
		qm.Load(orm.EnvironmentVariableRels.CreatedBy),
		qm.Load(orm.EnvironmentVariableRels.UpdatedBy),
	).All(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrEnvVarsNotFound(err, orgName, pipelineName, envName)
		default:
			return nil, err
		}
	}

	return envVars, nil
}

func loadEnvVarTx(ctx context.Context, tx *sql.Tx, orgName string, pipelineName string, envName string, varID int64) (*orm.Organization, *orm.Pipeline, *orm.Environment, *orm.EnvironmentVariable, *apperr.Error) {
	var appErr *apperr.Error
	org, pipeline, appErr := xorm.GetPpelineTx(ctx, tx, orgName, pipelineName)
	if appErr != nil {
		return nil, nil, nil, nil, appErr
	}

	env, appErr := xorm.GetOrgEnvironmentTx(ctx, tx, org, envName)
	if appErr != nil {
		return nil, nil, nil, nil, appErr
	}

	envVar, err := pipeline.EnvironmentVariables(
		orm.EnvironmentVariableWhere.ID.EQ(varID),
		orm.EnvironmentVariableWhere.EnvironmentID.EQ(env.ID),
		qm.OrderBy(orm.EnvironmentVariableColumns.Index+" DESC"),
		qm.Load(orm.EnvironmentVariableRels.CreatedBy),
		qm.Load(orm.EnvironmentVariableRels.UpdatedBy),
	).One(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, nil, nil, nil, ares.ErrEnvVarNotFound(err, orgName, pipelineName, envName, varID)
		default:
			return nil, nil, nil, nil, err
		}
	}

	return org, pipeline, env, envVar, nil
}

func getEnvironmentVariable(envVar *orm.EnvironmentVariable) *dto.GetEnvironmentVariableResponse {
	// Hide the values of secrets
	value := envVar.Value
	if envVar.Secret {
		value = ""
	}

	return &dto.GetEnvironmentVariableResponse{
		ID:        envVar.ID,
		Key:       envVar.Key,
		Value:     value,
		Secret:    envVar.Secret,
		CreatedAt: strfmt.DateTime(envVar.CreatedAt),
		CreatedBy: &dto.CreatedBy{
			ID:       envVar.R.CreatedBy.ID,
			FullName: envVar.R.CreatedBy.FullName,
		},
		UpdatedAt: strfmt.DateTime(envVar.UpdatedAt),
		UpdatedBy: &dto.UpdatedBy{
			ID:       envVar.R.UpdatedBy.ID,
			FullName: envVar.R.UpdatedBy.FullName,
		},
		Deployed: envVar.Deployed,
	}
}
