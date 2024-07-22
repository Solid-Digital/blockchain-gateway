package testhelper

import (
	"context"
	"database/sql"
	"time"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (h *Helper) HasActiveDeployments(pipeline *orm.Pipeline) bool {
	var org *orm.Organization
	var deployments []*orm.Deployment
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		org, err = orm.Organizations(orm.OrganizationWhere.ID.EQ(pipeline.OrganizationID)).One(ctx, tx)
		if err != nil {
			return err
		}

		deployments, err = orm.Deployments(orm.DeploymentWhere.PipelineID.EQ(pipeline.ID)).All(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})
	h.suite.Require().NoError(err)

	// This may actually not be true, but without a deployment record we don't have a fullName, so
	// it would be very hard to determine whether the pipeline has deployments on kube.
	if len(deployments) == 0 {
		return false
	}

	// It may be the case that polling 10 times is not enough...
	for _, deployment := range deployments {
		if h.HasActivePods(org.Name, deployment.FullName) {
			return true
		}
	}

	return false
}

func (h *Helper) HasActivePods(orgName, fullName string) bool {
	// It may be the case that polling 10 times is not enough...
	for i := 0; i < 10; i++ {
		_, err := h.ares.DeploymentService.GetDeploymentPods(orgName, fullName)
		if err != nil {
			return false
		}

		time.Sleep(time.Second)
	}

	return true
}
