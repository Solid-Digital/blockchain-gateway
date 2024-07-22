package factory

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/testhelper"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/volatiletech/sqlboiler/boil"
)

func (f *Factory) baseVersionForOrgAndUserAndBase(org *orm.Organization, user *orm.User, base *orm.Base, version string, public, create bool) *orm.BaseVersion {
	baseVersion := &orm.BaseVersion{
		BaseID:         base.ID,
		CreatedByID:    user.ID,
		UpdatedByID:    user.ID,
		Version:        version,
		DockerImageRef: "registry.unchain.io/unchainio/janus-v2:latest",
		Entrypoint:     "janus",
		Public:         public,
	}

	if !create {
		return baseVersion
	}

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return baseVersion.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return baseVersion
}

func (f *Factory) baseVersionForOrgAndUser(org *orm.Organization, user *orm.User, public, create bool) (*orm.BaseVersion, *orm.Base) {
	base := f.baseForOrgAndUser(org, user, public, create)

	return f.baseVersionForOrgAndUserAndBase(org, user, base, testhelper.RandomVersion(), public, create), base

}

func (f *Factory) BaseVersionOrgUser(public, create bool) (*orm.BaseVersion, *orm.Base, *orm.Organization, *orm.User) {
	org, user := f.OrganizationAndUser(create)

	baseVersion, base := f.baseVersionForOrgAndUser(org, user, public, create)

	return baseVersion, base, org, user
}

func (f *Factory) BaseVersionForOrgAndBase(public, create bool, org *orm.Organization, base *orm.Base, version string) *orm.BaseVersion {
	return f.baseVersionForOrgAndUserAndBase(org, f.UserFromOrg(org), base, version, public, create)
}

func (f *Factory) BaseVersionAndBase(public, create bool) (*orm.BaseVersion, *orm.Base) {
	baseVersion, base, _, _ := f.BaseVersionOrgUser(public, create)

	return baseVersion, base
}

func (f *Factory) BaseVersion(public, create bool) *orm.BaseVersion {
	baseVersion, _, _, _ := f.BaseVersionOrgUser(public, create)

	return baseVersion
}

func (f *Factory) BaseVersionAndBaseForOrg(org *orm.Organization, public, create bool) (*orm.BaseVersion, *orm.Base) {
	user := f.UserFromOrg(org)
	baseVersion, base := f.baseVersionForOrgAndUser(org, user, public, create)

	return baseVersion, base
}

func (f *Factory) BaseVersionForOrg(org *orm.Organization, public, create bool) *orm.BaseVersion {
	user := f.UserFromOrg(org)
	baseVersion, _ := f.baseVersionForOrgAndUser(org, user, public, create)

	return baseVersion
}
