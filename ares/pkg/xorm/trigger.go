package xorm

import (
	"context"
	"database/sql"
	stderr "errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"github.com/volatiletech/sqlboiler/boil"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func GetTriggerTx(ctx context.Context, tx *sql.Tx, orgName, name string, mods ...qm.QueryMod) (*orm.Trigger, *apperr.Error) {
	trigger, err := orm.Triggers(
		append(mods, orm.TriggerWhere.Name.EQ(name))...,
	).One(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrComponentNotFound(err, ares.ComponentTypeTrigger, orgName, name)
		default:
			return nil, err
		}
	}

	return trigger, nil
}

func GetAllTriggersTx(ctx context.Context, tx *sql.Tx, orgName string, mods ...qm.QueryMod) ([]*orm.Trigger, *apperr.Error) {
	triggers, err := orm.Triggers(
		append(
			mods,
		)...,
	).All(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrComponentsNotFound(err, ares.ComponentTypeTrigger, orgName)
		default:
			return nil, err
		}
	}

	return triggers, nil
}

func GetTriggerVersionTx(ctx context.Context, tx *sql.Tx, orgName string, trigger *orm.Trigger, version string, mods ...qm.QueryMod) (*orm.TriggerVersion, *apperr.Error) {
	triggerVersion, err := trigger.TriggerVersions(
		append(mods, orm.TriggerVersionWhere.Version.EQ(version))...,
	).One(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrComponentVersionNotFound(err, ares.ComponentTypeTrigger, orgName, trigger.Name, version)
		default:
			return nil, err
		}
	}

	return triggerVersion, nil
}

func CreateTriggerTx(ctx context.Context, tx *sql.Tx, trigger *orm.Trigger) *apperr.Error {
	err := trigger.Insert(ctx, tx, boil.Infer())
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.Conflict):
			return ares.ErrDuplicateComponent(err, ares.ComponentTypeTrigger, trigger.Name)
		default:
			return err
		}
	}

	return nil
}
