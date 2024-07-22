package testhelper

import (
	"bitbucket.org/unchain/ares/gen/orm"
	"context"
	"database/sql"
)

func (h *Helper) DBGetSubscriptionPlan(planID int64) *orm.SubscriptionPlan {
	var err error
	var planFromDB *orm.SubscriptionPlan

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		planFromDB, err = orm.SubscriptionPlans(orm.SubscriptionPlanWhere.ID.EQ(planID)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return planFromDB
}

func (h *Helper) DBGetSubscription(orgID int64) *orm.Subscription {
	var err error
	var sub *orm.Subscription

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		sub, err = orm.Subscriptions(orm.SubscriptionWhere.OrganizationID.EQ(orgID)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return sub
}
