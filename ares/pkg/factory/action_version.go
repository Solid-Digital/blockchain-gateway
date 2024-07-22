package factory

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper"
	"github.com/volatiletech/sqlboiler/boil"
)

func (f *Factory) actionVersionForOrgAndUserAndAction(org *orm.Organization, user *orm.User, action *orm.Action, version string, public bool, create bool) *orm.ActionVersion {
	actionVersion := &orm.ActionVersion{
		ActionID:     action.ID,
		CreatedByID:  user.ID,
		UpdatedByID:  user.ID,
		Version:      version,
		InputSchema:  f.IOSchemaJSON(),
		OutputSchema: f.IOSchemaJSON(),
		Public:       public,
	}

	if !create {
		return actionVersion
	}

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return actionVersion.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return actionVersion
}

func (f *Factory) actionVersionForOrgAndUser(org *orm.Organization, user *orm.User, public, create bool) (*orm.ActionVersion, *orm.Action) {
	action := f.actionForOrgAndUser(org, user, public, create)
	return f.actionVersionForOrgAndUserAndAction(org, user, action, testhelper.RandomVersion(), public, create), action
}

func (f *Factory) ActionVersionForOrgAndAction(public, create bool, org *orm.Organization, action *orm.Action, version string) *orm.ActionVersion {
	return f.actionVersionForOrgAndUserAndAction(org, f.UserFromOrg(org), action, version, public, create)
}

func (f *Factory) ActionVersionOrgUser(public, create bool) (*orm.ActionVersion, *orm.Action, *orm.Organization, *orm.User) {
	org, user := f.OrganizationAndUser(create)
	actionVersion, action := f.actionVersionForOrgAndUser(org, user, public, create)
	return actionVersion, action, org, user
}

func (f *Factory) ActionVersionAndAction(public, create bool) (*orm.ActionVersion, *orm.Action) {
	actionVersion, action, _, _ := f.ActionVersionOrgUser(public, create)

	return actionVersion, action
}

func (f *Factory) ActionVersion(public, create bool) *orm.ActionVersion {
	actionVersion, _, _, _ := f.ActionVersionOrgUser(public, create)

	return actionVersion
}

func (f *Factory) ActionVersionAndActionForOrg(org *orm.Organization, public, create bool) (*orm.ActionVersion, *orm.Action) {
	user := f.UserFromOrg(org)
	actionVersion, action := f.actionVersionForOrgAndUser(org, user, public, create)

	return actionVersion, action
}

func (f *Factory) ActionVersionForOrg(org *orm.Organization, public, create bool) *orm.ActionVersion {
	user := f.UserFromOrg(org)
	actionVersion, _ := f.actionVersionForOrgAndUser(org, user, public, create)

	return actionVersion
}
