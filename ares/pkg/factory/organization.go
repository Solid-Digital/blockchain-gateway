package factory

import (
	"context"
	"database/sql"
	"strings"

	"bitbucket.org/unchain/ares/pkg/testhelper"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/Pallinder/go-randomdata"
	"github.com/volatiletech/sqlboiler/boil"
)

func (f *Factory) Organization(create bool) *orm.Organization {
	user := f.DTOUser(create)
	displayName := testhelper.Randumb(randomdata.SillyName())

	org := &orm.Organization{
		CreatedByID: user.ID,
		UpdatedByID: user.ID,
		DisplayName: displayName,
		Name:        f.toOrganizationName(displayName),
	}

	if !create {
		return org
	}

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return org.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return org
}

func (f *Factory) OrganizationAndUser(create bool) (*orm.Organization, *orm.User) {
	org := f.Organization(create)
	user := f.User(create)

	if !create {
		return org, user
	}

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		err := org.AddUsers(ctx, tx, false, user)
		if err != nil {
			return err
		}

		return nil
	})

	f.suite.Require().NoError(err)

	return org, user
}

func (f *Factory) toOrganizationName(s string) string {
	s = strings.ToLower(s)
	strings.Replace(s, " ", "-", -1)

	return s
}
