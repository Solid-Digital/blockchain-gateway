package factory

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Pallinder/go-randomdata"

	"bitbucket.org/unchain/ares/pkg/testhelper"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (f *Factory) Environment(org *orm.Organization, user *orm.User) *orm.Environment {
	return f.EnvironmentWithName(testhelper.Randumb(fmt.Sprintf("env-%d", randomdata.Number(100))), org, user)
}

func (f *Factory) EnvironmentWithName(name string, org *orm.Organization, user *orm.User) *orm.Environment {
	env := &orm.Environment{
		Name:        name,
		CreatedByID: user.ID,
		UpdatedByID: user.ID,
		MaxReplicas: int64(5),
	}

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return org.AddEnvironments(ctx, tx, true, env)
	})

	f.suite.Require().NoError(err)

	return env
}
