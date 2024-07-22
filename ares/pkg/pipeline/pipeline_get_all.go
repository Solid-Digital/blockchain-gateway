package pipeline

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	v1 "k8s.io/api/apps/v1"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/dto"
)

func (s *Service) GetAllPipelines(orgName string) (dto.GetAllPipelinesResponse, *apperr.Error) {
	var ret dto.GetAllPipelinesResponse
	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		org, appErr := xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		pipelines, err := orm.Pipelines(
			qm.Load(orm.PipelineRels.Deployments),
			qm.Load(qm.Rels(orm.PipelineRels.Deployments, orm.DeploymentRels.Environment)),
			qm.Load(qm.Rels(orm.PipelineRels.Deployments, orm.DeploymentRels.Configuration)),
			qm.Load(qm.Rels(orm.PipelineRels.Deployments, orm.DeploymentRels.CreatedBy)),
			qm.Load(qm.Rels(orm.PipelineRels.Deployments, orm.DeploymentRels.UpdatedBy)),
			qm.Load(orm.PipelineRels.CreatedBy),
			qm.Load(orm.PipelineRels.UpdatedBy),
			orm.PipelineWhere.OrganizationID.EQ(org.ID)).All(ctx, tx)
		if err != nil {
			return ares.ParsePQErr(err)
		}

		envs, appErr := xorm.GetAllOrgEnvironmentsTx(ctx, tx, org)
		if appErr != nil {
			return appErr
		}

		ret = make(dto.GetAllPipelinesResponse, len(pipelines))

		for i, p := range pipelines {
			ret[i], appErr = s.toDTOPipeline(org.Name, p, envs)
			if appErr != nil {
				return appErr
			}
		}

		return nil
	})

	if appErr != nil {
		return nil, appErr
	}
	return ret, nil
}

func (s *Service) getKubeDeployment(orgName string, fullName string) *v1.Deployment {
	kd, err := s.service.kube.GetDeployment(orgName, fullName)
	if err != nil {
		s.log.Warnf("%+v", err)
		kd = &v1.Deployment{
			Status: v1.DeploymentStatus{
				AvailableReplicas: 0,
			},
		}
	}

	return kd
}
