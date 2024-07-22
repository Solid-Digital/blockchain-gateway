package testhelper

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (h *Helper) DBPublicBases() []*orm.Base {
	var publicBases []*orm.Base
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		publicBases, err = orm.Bases(orm.BaseWhere.Public.EQ(true)).All(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return publicBases
}

func (h *Helper) DBBaseExists(baseID int64) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.Bases(orm.BaseWhere.ID.EQ(baseID)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}

func (h *Helper) DBGetBase(baseID int64) *orm.Base {
	if !h.DBBaseExists(baseID) {
		return nil
	}

	var baseFromDB *orm.Base
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		baseFromDB, err = orm.Bases(orm.BaseWhere.ID.EQ(baseID)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return baseFromDB
}

func (h *Helper) DBBaseByNameExists(name string) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.Bases(orm.BaseWhere.Name.EQ(name)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}

func (h *Helper) DBGetBaseByName(name string) *orm.Base {
	if !h.DBBaseByNameExists(name) {
		return nil
	}

	var baseFromDB *orm.Base
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		baseFromDB, err = orm.Bases(orm.BaseWhere.Name.EQ(name)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return baseFromDB
}
