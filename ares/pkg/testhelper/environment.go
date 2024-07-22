package testhelper

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/volatiletech/sqlboiler/boil"
)

func (h *Helper) SetMaxReplicas(env *orm.Environment, replicas int64) {
	var err error
	env.MaxReplicas = replicas

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		_, err = env.Update(ctx, tx, boil.Infer())
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)
}
