package testhelper

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (h *Helper) DBDraftConfigurationExists(draftConfigurationID int64) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.DraftConfigurations(orm.DraftConfigurationWhere.ID.EQ(draftConfigurationID)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}

func (h *Helper) DBGetDraftConfiguration(draftConfigurationID int64) *orm.DraftConfiguration {
	if !h.DBDraftConfigurationExists(draftConfigurationID) {
		return nil
	}

	var draftConfigurationFromDB *orm.DraftConfiguration
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		draftConfigurationFromDB, err = orm.DraftConfigurations(
			qm.Load(orm.DraftConfigurationRels.BaseDraftConfiguration),
			qm.Load(orm.DraftConfigurationRels.TriggerDraftConfiguration),
			qm.Load(orm.DraftConfigurationRels.ActionDraftConfigurations),
			orm.DraftConfigurationWhere.ID.EQ(draftConfigurationID)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return draftConfigurationFromDB
}
