package xorm

import (
	"context"
	"database/sql"
	stderr "errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"github.com/volatiletech/sqlboiler/boil"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/ares"
)

func GetBaseTx(ctx context.Context, tx *sql.Tx, orgName, name string, mods ...qm.QueryMod) (*orm.Base, *apperr.Error) {
	base, err := orm.Bases(
		append(mods, orm.BaseWhere.Name.EQ(name))...,
	).One(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrComponentNotFound(err, ares.ComponentTypeBase, orgName, name)
		default:
			return nil, err
		}
	}

	return base, nil
}

func GetAllBasesTx(ctx context.Context, tx *sql.Tx, orgName string, mods ...qm.QueryMod) ([]*orm.Base, *apperr.Error) {
	bases, err := orm.Bases(
		append(
			mods,
		)...,
	).All(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrComponentsNotFound(err, ares.ComponentTypeBase, orgName)
		default:
			return nil, err
		}
	}

	return bases, nil
}

func GetBaseVersionTx(ctx context.Context, tx *sql.Tx, orgName string, base *orm.Base, version string, mods ...qm.QueryMod) (*orm.BaseVersion, *apperr.Error) {
	baseVersion, err := base.BaseVersions(
		append(mods, orm.BaseVersionWhere.Version.EQ(version))...,
	).One(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrComponentVersionNotFound(err, ares.ComponentTypeBase, orgName, base.Name, version)
		default:
			return nil, err
		}
	}

	return baseVersion, nil
}

func CreateBaseTx(ctx context.Context, tx *sql.Tx, base *orm.Base) *apperr.Error {
	err := base.Insert(ctx, tx, boil.Infer())
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.Conflict):
			return ares.ErrDuplicateComponent(err, ares.ComponentTypeBase, base.Name)
		default:
			return err
		}
	}

	return nil
}
