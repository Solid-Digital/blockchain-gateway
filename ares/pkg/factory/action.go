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

func (f *Factory) actionForOrgAndUser(org *orm.Organization, user *orm.User, public, create bool) *orm.Action {
	displayName := testhelper.Randumb(randomdata.SillyName())

	action := &orm.Action{
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
		return action
	}

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return action.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return action
}

func (f *Factory) ActionOrgUser(public, create bool) (*orm.Action, *orm.Organization, *orm.User) {
	org, user := f.OrganizationAndUser(create)

	return f.actionForOrgAndUser(org, user, public, create), org, user
}

func (f *Factory) Action(public, create bool) *orm.Action {
	action, _, _ := f.ActionOrgUser(public, create)

	return action
}

func (f *Factory) ActionForOrg(org *orm.Organization, public, create bool) *orm.Action {
	user := f.UserFromOrg(org)
	action := f.actionForOrgAndUser(org, user, public, create)

	return action
}
