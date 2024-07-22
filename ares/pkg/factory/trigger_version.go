package factory

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper"
	"github.com/volatiletech/sqlboiler/boil"
)

func (f *Factory) triggerVersionForOrgAndUserAndTrigger(org *orm.Organization, user *orm.User, trigger *orm.Trigger, version string, public, create bool) *orm.TriggerVersion {
	triggerVersion := &orm.TriggerVersion{
		TriggerID:    trigger.ID,
		CreatedByID:  user.ID,
		UpdatedByID:  user.ID,
		Version:      version,
		InputSchema:  f.IOSchemaJSON(),
		OutputSchema: f.IOSchemaJSON(),
		Public:       public,
	}

	if !create {
		return triggerVersion
	}

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return triggerVersion.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return triggerVersion
}

func (f *Factory) triggerVersionForOrgAndUser(org *orm.Organization, user *orm.User, public bool, create bool) (*orm.TriggerVersion, *orm.Trigger) {
	trigger := f.triggerForOrgAndUser(org, user, public, create)

	return f.triggerVersionForOrgAndUserAndTrigger(org, user, trigger, testhelper.RandomVersion(), public, create), trigger
}

func (f *Factory) triggerVersionForOrgAndUserWithFile(org *orm.Organization, user *orm.User, public bool, create bool, fileName string) (*orm.TriggerVersion, *orm.Trigger) {
	trigger := f.triggerForOrgAndUser(org, user, public, create)

	triggerVersion := &orm.TriggerVersion{
		TriggerID:   trigger.ID,
		Version:     testhelper.Randumb("alpha"),
		FileID:      fileName,
		FileName:    fileName,
		CreatedByID: user.ID,
		UpdatedByID: user.ID,
	}

	if !create {
		return triggerVersion, trigger
	}

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return triggerVersion.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return triggerVersion, trigger
}

func (f *Factory) TriggerVersionForOrgAndTrigger(public, create bool, org *orm.Organization, trigger *orm.Trigger, version string) *orm.TriggerVersion {
	return f.triggerVersionForOrgAndUserAndTrigger(org, f.UserFromOrg(org), trigger, version, public, create)
}

func (f *Factory) TriggerVersionOrgUser(public, create bool) (*orm.TriggerVersion, *orm.Trigger, *orm.Organization, *orm.User) {
	org, user := f.OrganizationAndUser(create)

	triggerVersion, trigger := f.triggerVersionForOrgAndUser(org, user, public, create)

	return triggerVersion, trigger, org, user
}

func (f *Factory) TriggerVersionAndTrigger(public, create bool) (*orm.TriggerVersion, *orm.Trigger) {
	triggerVersion, trigger, _, _ := f.TriggerVersionOrgUser(public, create)

	return triggerVersion, trigger
}

func (f *Factory) TriggerVersion(public bool, create bool) *orm.TriggerVersion {
	triggerVersion, _, _, _ := f.TriggerVersionOrgUser(public, create)

	return triggerVersion
}

func (f *Factory) TriggerVersionAndTriggerForOrg(org *orm.Organization, public, create bool) (*orm.TriggerVersion, *orm.Trigger) {
	user := f.UserFromOrg(org)
	triggerVersion, trigger := f.triggerVersionForOrgAndUser(org, user, public, create)

	return triggerVersion, trigger
}

func (f *Factory) TriggerVersionForOrg(org *orm.Organization, public, create bool) *orm.TriggerVersion {
	user := f.UserFromOrg(org)
	triggerVersion, _ := f.triggerVersionForOrgAndUser(org, user, public, create)

	return triggerVersion
}

func (f *Factory) TriggerVersionForOrgWithFile(org *orm.Organization, public, create bool, fileName string) *orm.TriggerVersion {
	triggerVersion, _ := f.triggerVersionForOrgAndUserWithFile(org, f.UserFromOrg(org), public, create, fileName)

	return triggerVersion
}
