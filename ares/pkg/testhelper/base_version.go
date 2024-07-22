package testhelper

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (h *Helper) DBBaseVersionExists(baseVersionID int64) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.BaseVersions(orm.BaseVersionWhere.ID.EQ(baseVersionID)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}

func (h *Helper) DBGetBaseVersion(baseVersionID int64) *orm.BaseVersion {
	if !h.DBBaseVersionExists(baseVersionID) {
		return nil
	}

	var baseVersionFromDB *orm.BaseVersion
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		baseVersionFromDB, err = orm.BaseVersions(orm.BaseVersionWhere.ID.EQ(baseVersionID)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return baseVersionFromDB
}
