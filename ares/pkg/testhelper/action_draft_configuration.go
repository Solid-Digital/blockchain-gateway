package testhelper

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (h *Helper) ActionDraftConfigurationExists(name string) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.ActionDraftConfigurations(orm.ActionDraftConfigurationWhere.Name.EQ(name)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}

func (h *Helper) GetActionDraftConfiguration(name string) *orm.ActionDraftConfiguration {
	if !h.ActionDraftConfigurationExists(name) {
		return nil
	}

	var actionDraftConfiguration *orm.ActionDraftConfiguration
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		actionDraftConfiguration, err = orm.ActionDraftConfigurations(orm.ActionDraftConfigurationWhere.Name.EQ(name)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return actionDraftConfiguration
}
