package testhelper

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (h *Helper) DBPublicActions() []*orm.Action {
	var publicAction []*orm.Action
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		publicAction, err = orm.Actions(orm.ActionWhere.Public.EQ(true)).All(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return publicAction
}

func (h *Helper) DBActionExists(actionID int64) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.Actions(orm.ActionWhere.ID.EQ(actionID)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}

func (h *Helper) DBGetAction(actionID int64) *orm.Action {
	if !h.DBActionExists(actionID) {
		return nil
	}

	var actionFromDB *orm.Action
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		actionFromDB, err = orm.Actions(orm.ActionWhere.ID.EQ(actionID)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return actionFromDB
}

func (h *Helper) DBActionByNameExists(name string) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.Actions(orm.ActionWhere.Name.EQ(name)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}

func (h *Helper) DBGetActionByName(name string) *orm.Action {
	if !h.DBActionByNameExists(name) {
		return nil
	}

	var actionFromDB *orm.Action
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		actionFromDB, err = orm.Actions(orm.ActionWhere.Name.EQ(name)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return actionFromDB
}
