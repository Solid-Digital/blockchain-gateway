package pipeline

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/volatiletech/sqlboiler/queries/qm"

	v1 "k8s.io/api/apps/v1"

	"github.com/go-openapi/strfmt"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) GetDeployment(orgName string, pipelineName string, envName string) (*dto.GetDeploymentResponse, *apperr.Error) {
	var pipeline *orm.Pipeline
	var env *orm.Environment
	var org *orm.Organization
	var deployment *orm.Deployment
	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		org, pipeline, appErr = xorm.GetPpelineTx(ctx, tx, orgName, pipelineName)
		if appErr != nil {
			return appErr
		}

		env, appErr := xorm.GetOrgEnvironmentTx(ctx, tx, org, envName)
		if appErr != nil {
			return appErr
		}

		deployment, appErr = xorm.GetDeploymentTx(ctx, tx, org, pipeline, env,
			qm.Load(orm.DeploymentRels.Configuration),
			qm.Load(orm.DeploymentRels.CreatedBy),
			qm.Load(orm.DeploymentRels.UpdatedBy),
		)
		if appErr != nil {
			return appErr
		}

		return nil
	})

	if appErr != nil {
		return nil, appErr
	}

	kd := s.getKubeDeployment(orgName, deployment.FullName)

	return getDeployment(pipeline, env, deployment, kd), nil
}

func getDeployment(pipeline *orm.Pipeline, env *orm.Environment, deployment *orm.Deployment, kubeDeployment *v1.Deployment) *dto.GetDeploymentResponse {
	return &dto.GetDeploymentResponse{
		ID:        &deployment.ID,
		CreatedAt: (*strfmt.DateTime)(&deployment.CreatedAt),
		CreatedBy: &dto.CreatedBy{
			ID:       deployment.R.CreatedBy.ID,
			FullName: deployment.R.CreatedBy.FullName,
		},
		UpdatedAt: (*strfmt.DateTime)(&deployment.UpdatedAt),
		UpdatedBy: &dto.UpdatedBy{
			ID:       deployment.R.UpdatedBy.ID,
			FullName: deployment.R.UpdatedBy.FullName,
		},
		DesiredReplicas:   &deployment.Replicas,
		AvailableReplicas: &[]int64{int64(kubeDeployment.Status.AvailableReplicas)}[0],
		Configuration: &dto.PipelineConfigurationRevision{
			CommitMessage: deployment.R.Configuration.CommitMessage,
			Revision:      deployment.R.Configuration.Revision,
		},
		Dirty: &deployment.Dirty,
		Host:  &deployment.Host,
		Image: &deployment.Image,
		Path:  &deployment.Path,
		URL:   &deployment.URL,
	}
}
