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

func (f *Factory) triggerForOrgAndUser(org *orm.Organization, user *orm.User, public bool, create bool) *orm.Trigger {
	displayName := testhelper.Randumb(randomdata.SillyName())

	trigger := &orm.Trigger{
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
		return trigger
	}

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return trigger.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return trigger
}

func (f *Factory) TriggerOrgUser(public bool, create bool) (*orm.Trigger, *orm.Organization, *orm.User) {
	org, user := f.OrganizationAndUser(create)

	return f.triggerForOrgAndUser(org, user, public, create), org, user
}

func (f *Factory) Trigger(public bool, create bool) *orm.Trigger {
	trigger, _, _ := f.TriggerOrgUser(public, create)

	return trigger
}

func (f *Factory) TriggerForOrg(org *orm.Organization, public bool, create bool) *orm.Trigger {
	return f.triggerForOrgAndUser(org, f.UserFromOrg(org), public, create)
}
