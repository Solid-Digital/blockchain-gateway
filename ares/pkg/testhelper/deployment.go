package testhelper

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (h *Helper) DBDeploymentExists(pipelineID, envID int64) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.Deployments(
			orm.DeploymentWhere.PipelineID.EQ(pipelineID),
			orm.DeploymentWhere.EnvironmentID.EQ(envID)).Exists(ctx, tx)

		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}

func (h *Helper) DBGetDeployment(pipelineID, envID int64) *orm.Deployment {
	if !h.DBDeploymentExists(pipelineID, envID) {
		return nil
	}

	var deployment *orm.Deployment
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		deployment, err = orm.Deployments(
			orm.DeploymentWhere.PipelineID.EQ(pipelineID),
			orm.DeploymentWhere.EnvironmentID.EQ(envID)).One(ctx, tx)

		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return deployment
}

func (h *Helper) DBDeploymentByIdExists(deploymentID int64) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.Deployments(
			orm.DeploymentWhere.ID.EQ(deploymentID)).Exists(ctx, tx)

		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}
