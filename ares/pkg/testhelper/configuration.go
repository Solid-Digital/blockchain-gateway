package testhelper

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (h *Helper) DBLatestConfigurationRevision(pipeline *orm.Pipeline) *orm.Configuration {
	var config *orm.Configuration

	err := h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		var err error

		config, err = pipeline.Configurations(
			qm.Load(orm.ConfigurationRels.CreatedBy),
			qm.Load(orm.ConfigurationRels.UpdatedBy),
			qm.Load(qm.Rels(orm.ConfigurationRels.BaseConfiguration, orm.BaseConfigurationRels.Version, orm.BaseVersionRels.Base)),
			qm.Load(qm.Rels(orm.ConfigurationRels.TriggerConfiguration, orm.TriggerConfigurationRels.Version, orm.TriggerVersionRels.Trigger)),
			qm.Load(qm.Rels(orm.ConfigurationRels.ActionConfigurations, orm.ActionConfigurationRels.Version, orm.ActionVersionRels.Action)),
			qm.OrderBy(orm.ConfigurationColumns.Revision+" DESC"),
		).One(ctx, tx)

		if err != nil {
			return err
		}

		return nil
	})
	h.suite.Require().NoError(err)

	return config
}
func (h *Helper) DBConfigurationExists(configurationID int64) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.Configurations(orm.ConfigurationWhere.ID.EQ(configurationID)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}

func (h *Helper) DBGetConfiguration(configurationID int64) *orm.Configuration {
	if !h.DBConfigurationExists(configurationID) {
		return nil
	}

	var configurationFromDB *orm.Configuration
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		configurationFromDB, err = orm.Configurations(
			qm.Load(orm.ConfigurationRels.BaseConfiguration),
			qm.Load(orm.ConfigurationRels.TriggerConfiguration),
			qm.Load(orm.ConfigurationRels.ActionConfigurations),
			orm.ConfigurationWhere.ID.EQ(configurationID)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return configurationFromDB
}
