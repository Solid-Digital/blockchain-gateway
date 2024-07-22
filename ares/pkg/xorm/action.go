package xorm

import (
	"context"
	"database/sql"
	stderr "errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/volatiletech/sqlboiler/boil"
)

func GetActionTx(ctx context.Context, tx *sql.Tx, orgName, name string, mods ...qm.QueryMod) (*orm.Action, *apperr.Error) {
	action, err := orm.Actions(
		append(mods, orm.ActionWhere.Name.EQ(name))...,
	).One(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrComponentNotFound(err, ares.ComponentTypeAction, orgName, name)
		default:
			return nil, err
		}
	}

	return action, nil
}

func GetAllActionsTx(ctx context.Context, tx *sql.Tx, orgName string, mods ...qm.QueryMod) ([]*orm.Action, *apperr.Error) {
	actions, err := orm.Actions(
		append(
			mods,
		)...,
	).All(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrComponentsNotFound(err, ares.ComponentTypeAction, orgName)
		default:
			return nil, err
		}
	}

	return actions, nil
}

func GetActionVersionTx(ctx context.Context, tx *sql.Tx, orgName string, action *orm.Action, version string, mods ...qm.QueryMod) (*orm.ActionVersion, *apperr.Error) {
	actionVersion, err := action.ActionVersions(
		append(mods, orm.ActionVersionWhere.Version.EQ(version))...,
	).One(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrComponentVersionNotFound(err, ares.ComponentTypeAction, orgName, action.Name, version)
		default:
			return nil, err
		}
	}

	return actionVersion, nil
}

func CreateActionTx(ctx context.Context, tx *sql.Tx, action *orm.Action) *apperr.Error {
	err := action.Insert(ctx, tx, boil.Infer())
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.Conflict):
			return ares.ErrDuplicateComponent(err, ares.ComponentTypeAction, action.Name)
		default:
			return err
		}
	}

	return nil
}
