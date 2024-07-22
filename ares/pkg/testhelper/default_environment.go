package testhelper

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (h *Helper) GetDefaultEnvironments() []*orm.DefaultEnvironment {
	var envs []*orm.DefaultEnvironment
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		envs, err = orm.DefaultEnvironments().All(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return envs
}

//Map default environments to map, using the name as key
func (h *Helper) DefaultEnvironmentsToMap(envs []*orm.DefaultEnvironment) map[string]*orm.DefaultEnvironment {
	m := map[string]*orm.DefaultEnvironment{}
	for _, env := range envs {
		m[env.Name] = env
	}

	return m
}
