package testhelper

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (h *Helper) DBPublicTriggers() []*orm.Trigger {
	var publicTriggers []*orm.Trigger
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		publicTriggers, err = orm.Triggers(orm.TriggerWhere.Public.EQ(true)).All(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return publicTriggers
}

func (h *Helper) DBTriggerExists(triggerID int64) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.Triggers(orm.TriggerWhere.ID.EQ(triggerID)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}

func (h *Helper) DBGetTrigger(triggerID int64) *orm.Trigger {
	if !h.DBTriggerExists(triggerID) {
		return nil
	}

	var triggerFromDB *orm.Trigger
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		triggerFromDB, err = orm.Triggers(orm.TriggerWhere.ID.EQ(triggerID)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return triggerFromDB
}

func (h *Helper) DBTriggerByNameExists(name string) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.Triggers(orm.TriggerWhere.Name.EQ(name)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}

func (h *Helper) DBGetTriggerByName(name string) *orm.Trigger {
	if !h.DBTriggerByNameExists(name) {
		return nil
	}

	var triggerFromDB *orm.Trigger
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		triggerFromDB, err = orm.Triggers(orm.TriggerWhere.Name.EQ(name)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return triggerFromDB
}
