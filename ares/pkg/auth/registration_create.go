package auth

import (
	"context"
	"database/sql"
	"strings"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"fmt"
	"time"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/volatiletech/null"
)

func (s *Service) CreateRegistration(params *dto.CreateRegistrationRequest) *apperr.Error {
	email := params.Email.String()
	orgName := *params.OrganizationName
	orgDisplayName := *params.OrganizationDisplayName

	// create sign-up token
	_, token, err := s.generateAuthCode()
	if err != nil {
		return apperr.Internal.Wrap(err)
	}

	var user *orm.User

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error
		var org *orm.Organization

		user, appErr = xorm.CreateUserTx(ctx, tx, email, token)
		if appErr != nil {
			return appErr
		}

		org, appErr = xorm.CreateOrganizationTx(ctx, tx, orgName, orgDisplayName, user)
		if appErr != nil {
			return appErr
		}

		// TODO: Move to xorm and refactor
		var provider *orm.OrganizationBillingProvider
		var subscription *orm.Subscription

		provider, subscription, appErr = s.createOrganizationSubscriptionTx(ctx, tx, params, user, org)
		if appErr != nil {
			return appErr
		}

		err = org.AddOrganizationBillingProviders(ctx, tx, true, provider)
		if err != nil {
			return ares.ParsePQErr(err)
		}

		err = org.AddSubscriptions(ctx, tx, true, subscription)
		if err != nil {
			return ares.ParsePQErr(err)
		}

		return nil
	})
	if appErr != nil {
		return appErr
	}

	err = s.enforcer.MakeUser(user.ID)
	if err != nil {
		return apperr.Internal.Wrap(err)
	}

	url := fmt.Sprintf(ares.URLFmt, s.cfg.ConnectURL, token, email, orgName, false)
	err = s.mailer.SendSignUpMessage(email, token, url)
	if err != nil {
		return apperr.Internal.Wrap(err)
	}

	return nil
}

// TODO: move to xorm and refactor
func (s *Service) createOrganizationSubscriptionTx(ctx context.Context, tx *sql.Tx, params *dto.CreateRegistrationRequest, user *orm.User, org *orm.Organization) (*orm.OrganizationBillingProvider, *orm.Subscription, *apperr.Error) {
	obp := &orm.OrganizationBillingProvider{
		OrganizationID: org.ID,
		CreatedByID:    user.ID,
		UpdatedByID:    user.ID,
	}

	subscription := &orm.Subscription{
		OrganizationID: org.ID,
		StartDate:      time.Now().UTC(),
		CreatedByID:    user.ID,
		UpdatedByID:    user.ID,
	}

	freePlan, err := orm.SubscriptionPlans(orm.SubscriptionPlanWhere.Name.EQ(ares.TierFreePlan)).One(ctx, tx)
	if err != nil {
		return nil, nil, ares.ParsePQErr(err)
	}

	switch strings.ToUpper(params.RegistrationReference.Source) {
	case ares.ProviderAWS:
		billingInfo, err := s.resolveAWSCustomer(params.RegistrationReference.Token)
		if err != nil {
			return nil, nil, apperr.Internal.Wrap(err)
		}
		obp.ProviderName = ares.ProviderAWS
		obp.BillingInfo = null.JSONFrom(billingInfo)

		subscription.Status = null.StringFrom(ares.StatusPendingConfirmation)
		subscription.SubscriptionPlanID = freePlan.ID
	default:
		obp.ProviderName = ares.ProviderUnchain
		obp.BillingInfo = null.JSONFrom([]byte(`{}`))

		subscription.SubscriptionPlanID = freePlan.ID
		subscription.Status = null.StringFrom(ares.StatusActive)
	}

	return obp, subscription, nil
}

func (s *Service) resolveAWSCustomer(awsToken string) (billingInfo []byte, err error) {
	// Call aws Marketplace Metering Resolve Customer API
	customerID, productCode, err := s.AWS.ResolveCustomer(awsToken)
	if err != nil {
		return nil, err
	}

	// store customerID in billing info JSON
	billingInfoString := fmt.Sprintf(`{"awsCustomerId": "%s", "productCode": "%s"}`, customerID, productCode)

	return []byte(billingInfoString), nil
}
