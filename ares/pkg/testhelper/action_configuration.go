package testhelper

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (h *Helper) ActionConfigurationExists(name string) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.ActionConfigurations(orm.ActionConfigurationWhere.Name.EQ(name)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}

func (h *Helper) GetActionConfiguration(name string) *orm.ActionConfiguration {
	if !h.ActionConfigurationExists(name) {
		return nil
	}

	var actionConfiguration *orm.ActionConfiguration
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		actionConfiguration, err = orm.ActionConfigurations(orm.ActionConfigurationWhere.Name.EQ(name)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return actionConfiguration
}
