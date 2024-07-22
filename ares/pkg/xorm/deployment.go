package xorm

import (
	"context"
	"database/sql"
	stderr "errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/ares"
)

func GetDeploymentTx(ctx context.Context, tx *sql.Tx, org *orm.Organization, pipeline *orm.Pipeline, env *orm.Environment, mods ...qm.QueryMod) (*orm.Deployment, *apperr.Error) {
	deployment, err := pipeline.Deployments(
		append(mods, orm.DeploymentWhere.EnvironmentID.EQ(env.ID))...,
	).One(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrDeploymentNotFound(err, org.Name, pipeline.Name, env.Name)
		default:
			return nil, err
		}
	}

	return deployment, nil
}
