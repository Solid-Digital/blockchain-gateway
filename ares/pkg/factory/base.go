package factory

import (
	"context"
	"database/sql"
	"time"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper"
	"github.com/Pallinder/go-randomdata"
	"github.com/volatiletech/sqlboiler/boil"
)

func (f *Factory) baseForOrgAndUser(org *orm.Organization, user *orm.User, public, create bool) *orm.Base {
	displayName := testhelper.Randumb(randomdata.SillyName())

	base := &orm.Base{
		DeveloperID: org.ID,
		CreatedAt:   time.Time{},
		CreatedByID: user.ID,
		UpdatedAt:   time.Time{},
		UpdatedByID: user.ID,
		Name:        f.toOrganizationName(displayName),
		DisplayName: displayName,
		Public:      public,
	}

	if !create {
		return base
	}

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return base.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return base
}

func (f *Factory) BaseOrgUser(public, create bool) (*orm.Base, *orm.Organization, *orm.User) {
	org, user := f.OrganizationAndUser(create)

	return f.baseForOrgAndUser(org, user, public, create), org, user
}

func (f *Factory) Base(public, create bool) *orm.Base {
	base, _, _ := f.BaseOrgUser(public, create)

	return base
}

func (f *Factory) BaseForOrg(org *orm.Organization, public, create bool) *orm.Base {
	return f.baseForOrgAndUser(org, f.UserFromOrg(org), public, create)
}
