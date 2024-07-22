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

func GetOrgEnvironmentTx(ctx context.Context, tx *sql.Tx, org *orm.Organization, envName string) (*orm.Environment, *apperr.Error) {
	env, err := org.Environments(orm.EnvironmentWhere.Name.EQ(envName)).One(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrEnvironmentNotFound(err, org.Name, envName)
		default:
			return nil, err
		}
	}

	return env, nil
}

func GetAllOrgEnvironmentsTx(ctx context.Context, tx *sql.Tx, org *orm.Organization) ([]*orm.Environment, *apperr.Error) {
	envs, err := org.Environments(qm.OrderBy(orm.EnvironmentColumns.Index)).All(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrEnvironmentsNotFound(err, org.Name)
		default:
			return nil, err
		}
	}

	return envs, nil
}
