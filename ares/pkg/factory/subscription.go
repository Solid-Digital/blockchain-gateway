package factory

import (
	"bitbucket.org/unchain/ares/gen/orm"
	"context"
	"database/sql"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"time"
)

func (f *Factory) SubscriptionPlan(create bool, planName string) *orm.SubscriptionPlan {
	user := f.DTOUser(create)

	plan := &orm.SubscriptionPlan{
		Name:          planName,
		PipelineLimit: 3,
		CreatedByID:   user.ID,
		UpdatedByID:   user.ID,
	}

	if !create {
		return plan
	}

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return plan.Upsert(ctx, tx, false, []string{}, boil.Infer(), boil.Infer())
	})

	f.suite.Require().NoError(err)

	return plan
}

func (f *Factory) Subscription(create bool, orgID int64, planName, status string) *orm.Subscription {
	user := f.User(create)
	subs := &orm.Subscription{
		SubscriptionPlanID: 1,
		OrganizationID:     orgID,
		Status:             null.StringFrom(status),
		StartDate:          time.Now().UTC(),
		CreatedByID:        user.ID,
		UpdatedByID:        user.ID,
	}

	if !create {
		return subs
	}

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return subs.Upsert(ctx, tx, false, []string{}, boil.Infer(), boil.Infer())
	})

	f.suite.Require().NoError(err)

	return subs
}
