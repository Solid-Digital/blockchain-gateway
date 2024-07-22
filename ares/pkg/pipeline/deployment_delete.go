package pipeline

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"
)

func (s *Service) RemoveDeployment(orgName string, pipelineName string, envName string) *apperr.Error {
	return ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
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

		deployment, appErr := xorm.GetDeploymentTx(ctx, tx, org, pipeline, env)
		if appErr != nil {
			return appErr
		}

		err = s.service.kube.DeleteDeployment(orgName, deployment.FullName)
		if err != nil {
			return apperr.Internal.Wrap(err)
		}

		_, err = deployment.Delete(ctx, tx)
		if err != nil {
			return ares.ParsePQErr(err)
		}

		return nil
	})
}
