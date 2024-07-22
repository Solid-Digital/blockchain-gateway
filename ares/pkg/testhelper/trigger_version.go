package testhelper

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (h *Helper) DBTriggerVersionExists(triggerVersionID int64) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.TriggerVersions(orm.TriggerVersionWhere.ID.EQ(triggerVersionID)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}

func (h *Helper) DBGetTriggerVersion(triggerVersionID int64) *orm.TriggerVersion {
	if !h.DBTriggerVersionExists(triggerVersionID) {
		return nil
	}

	var triggerVersionFromDB *orm.TriggerVersion
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		triggerVersionFromDB, err = orm.TriggerVersions(orm.TriggerVersionWhere.ID.EQ(triggerVersionID)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return triggerVersionFromDB
}
