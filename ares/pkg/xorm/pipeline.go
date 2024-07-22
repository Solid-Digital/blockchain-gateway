package xorm

import (
	"context"
	"database/sql"
	stderr "errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/ares"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

func GetPpelineTx(ctx context.Context, tx *sql.Tx, orgName string, pipelineName string, mods ...qm.QueryMod) (*orm.Organization, *orm.Pipeline, *apperr.Error) {
	var appErr *apperr.Error

	org, appErr := GetOrganizationTx(ctx, tx, orgName)
	if appErr != nil {
		return nil, nil, appErr
	}

	mods = append(
		mods,

		orm.PipelineWhere.OrganizationID.EQ(org.ID),
		orm.PipelineWhere.Name.EQ(pipelineName),
	)

	pipeline, err := orm.Pipelines(
		mods...,
	).One(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, nil, ares.ErrPipelineNotFound(err, orgName, pipelineName)
		default:
			return nil, nil, err
		}
	}

	return org, pipeline, nil
}
