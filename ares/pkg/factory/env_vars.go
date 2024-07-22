package factory

import (
	"context"
	"database/sql"
	"math/rand"
	"sort"

	"github.com/volatiletech/sqlboiler/boil"

	"bitbucket.org/unchain/ares/pkg/testhelper"
	"github.com/Pallinder/go-randomdata"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (f *Factory) MakeRandomVars(org *orm.Organization, pipeline *orm.Pipeline, env *orm.Environment, user *orm.User, params EnvVarFactoryParams) orm.EnvironmentVariableSlice {
	n := randomdata.Number(5, 100)

	vars := make([]*orm.EnvironmentVariable, n)
	indices := rand.Perm(n)

	for i := 0; i < n; i++ {
		var secret bool

		switch params {
		case OnlySecrets:
			secret = true
		case OnlyVars:
			secret = false
		case BothSecretsAndVars:
			fallthrough
		default:
			secret = randomdata.Boolean()
		}
		vars[i] = &orm.EnvironmentVariable{
			OrganizationID: org.ID,
			PipelineID:     pipeline.ID,
			EnvironmentID:  env.ID,
			CreatedByID:    user.ID,
			UpdatedByID:    user.ID,

			Index:    int64(indices[i]),
			Key:      testhelper.Randumb(randomdata.Noun()),
			Value:    testhelper.Randumb(randomdata.Noun()),
			Secret:   secret,
			Deployed: false,
		}
	}

	return vars
}

type EnvVarFactoryParams int

const (
	OnlySecrets EnvVarFactoryParams = iota
	OnlyVars
	BothSecretsAndVars
)

func (f *Factory) EnvVars(create bool, user *orm.User, org *orm.Organization, pipeline *orm.Pipeline, env *orm.Environment, params EnvVarFactoryParams) (vars []*orm.EnvironmentVariable) {
	vars = f.MakeRandomVars(org, pipeline, env, user, params)

	if !create {
		return vars
	}

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		for _, v := range vars {
			err := v.Insert(ctx, tx, boil.Infer())
			if err != nil {
				return err
			}
		}

		return nil
	})
	f.suite.Require().NoError(err)

	sort.Slice(vars, func(i, j int) bool {
		return vars[i].Index > vars[j].Index
	})

	return vars
}
