package testhelper

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (h *Helper) DBOrgExists(orgName string) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.Organizations(orm.OrganizationWhere.Name.EQ(orgName)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	xrequire.NoError(h.suite.T(), err)

	return exists
}

func (h *Helper) DBOrganizationCount() int64 {
	var count int64
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		count, err = orm.Organizations().Count(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	xrequire.NoError(h.suite.T(), err)

	return count
}

func (h *Helper) DBGetOrgByName(orgName string) *orm.Organization {
	var orgFromDB *orm.Organization
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		orgFromDB, err = orm.Organizations(
			orm.OrganizationWhere.Name.EQ(orgName),
			qm.Load(orm.OrganizationRels.OrganizationBillingProviders),
			qm.Load(orm.OrganizationRels.Subscriptions),
			qm.Load(orm.OrganizationRels.Environments)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	xrequire.NoError(h.suite.T(), err)

	return orgFromDB
}

func (h *Helper) DBGetOrgByID(orgID int64) *orm.Organization {
	var orgFromDB *orm.Organization
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		orgFromDB, err = orm.Organizations(orm.OrganizationWhere.ID.EQ(orgID),
			qm.Load(orm.OrganizationRels.OrganizationBillingProviders),
			qm.Load(orm.OrganizationRels.Subscriptions),
			qm.Load(orm.OrganizationRels.Environments)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	xrequire.NoError(h.suite.T(), err)

	return orgFromDB
}
