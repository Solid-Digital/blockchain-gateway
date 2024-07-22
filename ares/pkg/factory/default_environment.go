package factory

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper"
	"github.com/Pallinder/go-randomdata"
	"github.com/volatiletech/sqlboiler/boil"
)

func (f *Factory) DefaultEnvironment(create bool) *orm.DefaultEnvironment {
	env := &orm.DefaultEnvironment{
		Name:        testhelper.Randumb(randomdata.SillyName()),
		MaxReplicas: int64(randomdata.Number(9)),
	}

	if !create {
		return env
	}

	f.suite.Require().NotNil(f.ares)

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return env.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return env
}
