package pipeline

import (
	"context"
	"database/sql"
	"fmt"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/ares"

	"github.com/volatiletech/sqlboiler/boil"

	"bitbucket.org/unchain/ares/gen/dto"
)

func (s *Service) UpdateEnvironmentVariable(params *dto.UpdateEnvironmentVariablesRequest, orgName string, pipelineName string, envName string, varID int64, user *dto.User) (*dto.GetEnvironmentVariableResponse, *apperr.Error) {
	var ret *dto.GetEnvironmentVariableResponse
	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var err error

		org, pipeline, env, envVar, appErr := loadEnvVarTx(ctx, tx, orgName, pipelineName, envName, varID)
		if appErr != nil {
			return appErr
		}

		if params.Key != nil {
			envVar.Key = *params.Key
		}

		if params.Value != nil {
			envVar.Value = *params.Value
		}

		envVar.UpdatedByID = user.ID
		envVar.Deployed = false

		_, err = envVar.Update(ctx, tx, boil.Infer())
		if err != nil {
			return ares.ParsePQErr(err)
		}

		appErr = markExistingDeploymentDirtyTx(ctx, tx, true, org, pipeline, env)
		if appErr != nil {
			return appErr
		}

		err = envVar.L.LoadUpdatedBy(ctx, tx, true, envVar, nil)
		if err != nil {
			return ares.ErrLoadingUpdatedBy(ares.ParsePQErr(err), fmt.Sprintf("%d", varID))
		}

		ret = getEnvironmentVariable(envVar)

		return nil
	})

	if appErr != nil {
		return nil, appErr
	}

	return ret, nil
}
