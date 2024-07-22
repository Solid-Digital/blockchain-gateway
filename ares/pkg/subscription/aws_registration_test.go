package subscription_test

import (
	"testing"
	"time"

	mock_ares "bitbucket.org/unchain/ares/gen/mocks"
	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/subscription"
	"bitbucket.org/unchain/ares/pkg/testhelper"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
	"github.com/Pallinder/go-randomdata"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null"
)

func (s *TestSuite) TestService_HandleNotificationMessage() {
	productCode := testhelper.Randumb(randomdata.Noun())

	customerID1 := testhelper.Randumb(randomdata.Noun())
	customerID2 := testhelper.Randumb(randomdata.Noun())
	customerID3 := testhelper.Randumb(randomdata.Noun())

	_ = s.factory.SubscriptionPlan(true, ares.TierFreePlan)
	_ = s.factory.SubscriptionPlan(true, ares.TierStarterPlan)

	org1, _ := s.factory.OrganizationAndAWSBillingProvider(true, productCode, customerID1)
	org2, _ := s.factory.OrganizationAndAWSBillingProvider(true, productCode, customerID2)
	org3, _ := s.factory.OrganizationAndAWSBillingProvider(true, productCode, customerID3)

	mpm1 := ares.AWSMarketplaceNotificationMessage{
		Body: ares.AWSMarketplaceNotificationMessageBody{
			CustomerIdentifier: customerID1,
			ProductCode:        productCode,
			Action:             subscription.ActionSubscribeSuccess,
		},
		Handle: s.helper.StringPtr(testhelper.Randumb(randomdata.Noun())),
	}
	mpm2 := ares.AWSMarketplaceNotificationMessage{
		Body: ares.AWSMarketplaceNotificationMessageBody{
			CustomerIdentifier: customerID2,
			ProductCode:        productCode,
			Action:             subscription.ActionEntitlementUpdated,
		},
		Handle: s.helper.StringPtr(testhelper.Randumb(randomdata.Noun())),
	}
	mpm3 := ares.AWSMarketplaceNotificationMessage{
		Body: ares.AWSMarketplaceNotificationMessageBody{
			CustomerIdentifier: customerID3,
			ProductCode:        productCode,
			Action:             subscription.ActionUnsubscribeSuccess,
		},
		Handle: s.helper.StringPtr(testhelper.Randumb(randomdata.Noun())),
	}

	sub1 := s.factory.Subscription(true, org1.ID, ares.TierFreePlan, ares.StatusInactive)
	sub2 := s.factory.Subscription(true, org2.ID, ares.TierFreePlan, ares.StatusInactive)
	sub3 := s.factory.Subscription(true, org3.ID, ares.TierFreePlan, ares.StatusInactive)

	inOneMonth := time.Now().Add(24 * 365 / 12 * time.Hour)
	past := time.Date(2019, 6, 21, 12, 0, 0, 0, time.UTC)
	entitlements1 := []*ares.AWSEntitlement{{
		ProductCode:        testhelper.StringPtr(productCode),
		CustomerIdentifier: testhelper.StringPtr(customerID1),
		Dimension:          testhelper.StringPtr(subscription.EntitlementFree),
		ExpirationDate:     &inOneMonth,
	}}
	entitlements2 := []*ares.AWSEntitlement{{
		ProductCode:        testhelper.StringPtr(productCode),
		CustomerIdentifier: testhelper.StringPtr(customerID2),
		Dimension:          testhelper.StringPtr(subscription.EntitlementStarter),
		ExpirationDate:     &inOneMonth,
	}}
	entitlements3 := []*ares.AWSEntitlement{{
		ProductCode:        testhelper.StringPtr(productCode),
		CustomerIdentifier: testhelper.StringPtr(customerID3),
		Dimension:          testhelper.StringPtr(subscription.EntitlementFree),
		ExpirationDate:     &past,
	}}

	mock := mock_ares.NewMockAWSClient(gomock.NewController(s.T()))
	mock.EXPECT().DeleteSQSMessage(gomock.Any()).AnyTimes().MinTimes(1)

	cases := map[string]struct {
		Message                    ares.AWSMarketplaceNotificationMessage
		Entitlements               []*ares.AWSEntitlement
		Organization               *orm.Organization
		Subscription               *orm.Subscription
		ExpectedSubscriptionStatus string
		ExpectedPlan               string
		// Valid means that start date < time.now > end date
		ExpectedValid bool
		Success       bool
	}{
		"subscribe success message": {
			mpm1,
			entitlements1,
			org1,
			sub1,
			ares.StatusActive,
			ares.TierFreePlan,
			true,
			true,
		},
		"entitlement update message": {
			mpm2,
			entitlements2,
			org2,
			sub2,
			ares.StatusActive,
			ares.TierStarterPlan,
			true,
			true,
		},
		"unsubscribe success message": {
			mpm3,
			entitlements3,
			org3,
			sub3,
			ares.StatusInactive,
			ares.TierFreePlan,
			false,
			true,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			mock.EXPECT().ReceiveMarketplaceNotification().Return(tc.Message)
			mock.EXPECT().GetEntitlements(gomock.Any(), gomock.Any()).Return(tc.Entitlements, nil)
			var err error
			cfg := &subscription.Config{
				ActivateHandler: false,
			}
			s.service, err = subscription.NewService(cfg, s.ares.DB, mock, s.ares.Log)
			xrequire.NoError(t, err)

			err = s.service.ConsumeMarketplaceNotificationMessage()

			if tc.Success {
				xrequire.NoError(t, err)

				sub := s.helper.DBGetSubscription(tc.Organization.ID)
				spFinal := s.helper.DBGetSubscriptionPlan(sub.SubscriptionPlanID)

				// then

				// Subscription should be active
				require.Equal(t, null.StringFrom(tc.ExpectedSubscriptionStatus), sub.Status)
				// Subscription should be starter
				require.Equal(t, tc.ExpectedPlan, spFinal.Name, "for org: %v", tc.Organization.ID)
				if tc.ExpectedValid {
					// Subscription start date is updated
					require.True(t, tc.Subscription.StartDate.Before(sub.StartDate))
					require.True(t, sub.StartDate.Before(time.Now()))
					require.True(t, sub.EndDate.After(time.Now()))
				} else {
					require.True(t, time.Now().Sub(sub.EndDate) > 0)
				}

			} else {
				xrequire.Error(t, err)
			}
		})
	}
}

// TODO Test expiration date and auto renewal
