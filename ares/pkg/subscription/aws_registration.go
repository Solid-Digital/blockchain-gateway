package subscription

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/unchainio/pkg/errors"
	"github.com/unchainio/pkg/iferr"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

const (
	ActionSubscribeSuccess   = "subscribe-success"
	ActionSubscribeFail      = "subscribe-fail"
	ActionUnsubscribePending = "unsubscribe-pending"
	ActionUnsubscribeSuccess = "unsubscribe-success"
	ActionEntitlementUpdated = "entitlement-updated"
	EntitlementFree          = "free"
	EntitlementStarter       = "starter"
	EntitlementBusiness      = "business"
	EntitlementEnterprise    = "enterprise"
)

func (s *Service) getSubscriptionPlanTx(ctx context.Context, tx *sql.Tx, customerID, productCode string) (planID int64, endDate *time.Time, err error) {
	entitlements, err := s.AWS.GetEntitlements(customerID, productCode)
	if err != nil {
		return 0, nil, err
	}

	var plan *orm.SubscriptionPlan
	for _, e := range entitlements {
		if e.Dimension == nil {
			continue
		}

		var tier string
		switch *e.Dimension {
		case EntitlementFree:
			tier = ares.TierFreePlan
		case EntitlementStarter:
			tier = ares.TierStarterPlan
		case EntitlementBusiness:
			tier = ares.TierBusinessPlan
		case EntitlementEnterprise:
			tier = ares.TierEnterprisePlan
		default:
			continue
		}
		plan, err = orm.SubscriptionPlans(orm.SubscriptionPlanWhere.Name.EQ(tier)).One(ctx, tx)
		if err != nil {
			return 0, nil, err
		} else {
			endDate = e.ExpirationDate
			break
		}
	}
	return plan.ID, endDate, nil

}

func (s *Service) handleAwsMarketplaceNotifications() {
	go func() {
		for {
			err := s.ConsumeMarketplaceNotificationMessage()
			iferr.Warn(err)
		}
	}()
}

func (s *Service) ConsumeMarketplaceNotificationMessage() error {
	msg := s.AWS.ReceiveMarketplaceNotification()

	if msg.Error != nil {
		iferr.Warn(errors.WithMessagef(msg.Error, "message handle: %+v", msg.Handle))
		return msg.Error
	}

	switch msg.Body.Action {
	case ActionSubscribeSuccess:
		err := s.handleAwsActionSubscribeSuccess(msg)
		if err != nil {
			return err
		}
	case ActionUnsubscribePending:
		err := s.handleAwsActionUnsubscribePending(msg)
		if err != nil {
			return err
		}
	case ActionSubscribeFail:
		err := s.handleAwsSubscribeFail(msg)
		if err != nil {
			return err
		}
	case ActionUnsubscribeSuccess:
		err := s.handleAwsUnsubscribeSuccess(msg)
		if err != nil {
			return err
		}
	case ActionEntitlementUpdated:
		err := s.handleAwsEntitlementUpdated(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) handleAwsActionSubscribeSuccess(msg ares.AWSMarketplaceNotificationMessage) error {
	s.log.Printf("#####\nSubscribe success message for product %s, customer identifier: \n%v\n#####\n", msg.Body.ProductCode, msg.Body.CustomerIdentifier)
	err := s.db.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		query := fmt.Sprintf("SELECT * FROM %s WHERE %s->>'%s' = $1", orm.TableNames.OrganizationBillingProvider, orm.OrganizationBillingProviderColumns.BillingInfo, "awsCustomerId")
		obp, err := orm.OrganizationBillingProviders(qm.SQL(query, msg.Body.CustomerIdentifier), qm.Load(orm.OrganizationBillingProviderRels.Organization)).One(ctx, tx)
		if err != nil {
			return err
		}

		var billingInfo ares.AWSBillingInfo
		err = json.Unmarshal(obp.BillingInfo.JSON, &billingInfo)
		if err != nil {
			return errors.Wrap(err, "billingInfo")
		}

		planID, expirationDate, err := s.getSubscriptionPlanTx(ctx, tx, billingInfo.CustomerIdentifier, billingInfo.ProductCode)
		if err != nil {
			return err
		}

		sub, err := obp.R.Organization.Subscriptions().One(ctx, tx)
		if err != nil {
			return err
		}

		sub.SubscriptionPlanID = planID
		sub.Status = null.StringFrom(ares.StatusActive)
		sub.StartDate = time.Now().UTC()
		if expirationDate != nil {
			sub.EndDate = *expirationDate
		}

		_, err = sub.Update(ctx, tx, boil.Infer())
		if err != nil {
			return err
		}

		err = s.AWS.DeleteSQSMessage(msg.Handle)
		if err != nil {
			return errors.Wrap(err, "error removing SQS message")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) handleAwsSubscribeFail(msg ares.AWSMarketplaceNotificationMessage) error {
	s.log.Printf("#####\nSubscribe fail message: \n%v\n#####\n", msg)
	return nil
}

func (s *Service) handleAwsActionUnsubscribePending(msg ares.AWSMarketplaceNotificationMessage) error {
	s.log.Printf("#####\nUnsubscribe pending message: \n%v\n#####\n", msg)
	return nil
}

func (s *Service) handleAwsUnsubscribeSuccess(msg ares.AWSMarketplaceNotificationMessage) error {
	s.log.Printf("#####\nUnsubscribe sucess message: \n%v\n#####\n", msg)
	err := s.db.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		query := fmt.Sprintf("SELECT * FROM %s WHERE %s->>'%s' = $1", orm.TableNames.OrganizationBillingProvider, orm.OrganizationBillingProviderColumns.BillingInfo, "awsCustomerId")
		obp, err := orm.OrganizationBillingProviders(qm.SQL(query, msg.Body.CustomerIdentifier), qm.Load(orm.OrganizationBillingProviderRels.Organization)).One(ctx, tx)
		if err != nil {
			return err
		}

		var billingInfo ares.AWSBillingInfo
		err = json.Unmarshal(obp.BillingInfo.JSON, &billingInfo)
		if err != nil {
			return errors.Wrap(err, "billingInfo")
		}

		_, expirationDate, err := s.getSubscriptionPlanTx(ctx, tx, billingInfo.CustomerIdentifier, billingInfo.ProductCode)
		if err != nil {
			return err
		}

		sub, err := obp.R.Organization.Subscriptions().One(ctx, tx)
		if err != nil {
			return err
		}

		sub.Status = null.StringFrom(ares.StatusInactive)
		if expirationDate != nil {
			sub.EndDate = *expirationDate
		}

		_, err = sub.Update(ctx, tx, boil.Infer())
		if err != nil {
			return err
		}

		// Placeholder: if applicable within 1 hr send metering data to AWS

		err = s.AWS.DeleteSQSMessage(msg.Handle)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) handleAwsEntitlementUpdated(msg ares.AWSMarketplaceNotificationMessage) error {
	s.log.Printf("######\nUpdate entitlement message: \n%v\n#####\n", msg)
	err := s.db.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		query := fmt.Sprintf("SELECT * FROM %s WHERE %s->>'%s' = $1", orm.TableNames.OrganizationBillingProvider, orm.OrganizationBillingProviderColumns.BillingInfo, "awsCustomerId")
		obp, err := orm.OrganizationBillingProviders(qm.SQL(query, msg.Body.CustomerIdentifier), qm.Load(orm.OrganizationBillingProviderRels.Organization)).One(ctx, tx)
		if err != nil {
			return err
		}

		var billingInfo ares.AWSBillingInfo
		err = json.Unmarshal(obp.BillingInfo.JSON, &billingInfo)
		if err != nil {
			return errors.Wrap(err, "billingInfo")
		}

		// FIXME in tests it will not get the correct planID based on the customer identifier and productCode
		planID, endDate, err := s.getSubscriptionPlanTx(ctx, tx, billingInfo.CustomerIdentifier, billingInfo.ProductCode)
		if err != nil {
			return err
		}

		sub, err := obp.R.Organization.Subscriptions().One(ctx, tx)
		if err != nil {
			return err
		}

		if planID != sub.SubscriptionPlanID {
			sub.SubscriptionPlanID = planID
			sub.Status = null.StringFrom(ares.StatusActive)
			sub.StartDate = time.Now().UTC()
			if endDate != nil {
				sub.EndDate = *endDate
			}
		}

		_, err = sub.Update(ctx, tx, boil.Infer())
		if err != nil {
			return err
		}

		err = s.AWS.DeleteSQSMessage(msg.Handle)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
