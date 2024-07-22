package pipeline

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/gen/dto"
)

func (s *Service) DeleteEnvironmentVariable(orgName string, pipelineName string, envName string, varID int64, user *dto.User) *apperr.Error {
	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var err error
		var appErr *apperr.Error

		org, pipeline, env, envVar, appErr := loadEnvVarTx(ctx, tx, orgName, pipelineName, envName, varID)
		if appErr != nil {
			return appErr
		}

		_, err = envVar.Delete(ctx, tx)
		if err != nil {
			return ares.ParsePQErr(err)
		}

		appErr = markExistingDeploymentDirtyTx(ctx, tx, true, org, pipeline, env)
		if appErr != nil {
			return appErr
		}

		return nil
	})

	if appErr != nil {
		return appErr
	}

	return nil
}
