package pipeline

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) DeletePipeline(orgName string, pipelineName string, principal *dto.User) *apperr.Error {
	var pipeline *orm.Pipeline
	var deployments []*orm.Deployment

	return ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error
		var err error

		_, pipeline, appErr = xorm.GetPpelineTx(ctx, tx, orgName, pipelineName)
		if appErr != nil {
			return appErr
		}

		deployments, err = orm.Deployments(orm.DeploymentWhere.PipelineID.EQ(pipeline.ID)).All(ctx, tx)
		if err != nil {
			return ares.ParsePQErr(err)
		}

		if len(deployments) > 0 {
			var names []string
			for _, deployment := range deployments {
				names = append(names, deployment.FullName)
			}

			// If this fails, possible removed pods, ingresses and env vars won't roll back
			err = s.service.kube.DeleteAllDeployments(orgName, names)
			if err != nil {
				return apperr.Internal.Wrap(err)
			}
		}

		// This will also delete configurations and deployments
		_, err = pipeline.Delete(ctx, tx)
		if err != nil {
			return ares.ParsePQErr(err)
		}

		return nil
	})
}
