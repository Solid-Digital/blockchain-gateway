package testhelper

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (h *Helper) DBEnvVarExists(envVarID int64) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.EnvironmentVariables(orm.EnvironmentVariableWhere.ID.EQ(envVarID)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}
