package auth_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/ares"

	"github.com/go-openapi/strfmt"

	"bitbucket.org/unchain/ares/pkg/testhelper"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
	"github.com/Pallinder/go-randomdata"

	"github.com/golang/mock/gomock"

	mock_ares "bitbucket.org/unchain/ares/gen/mocks"

	"github.com/stretchr/testify/require"

	"time"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/pkg/auth"
	"github.com/volatiletech/null"
)

func (s *TestSuite) TestService_CreateRegistration() {
	_ = s.factory.SubscriptionPlan(true, ares.TierFreePlan)

	// given
	registrationReference := &dto.CreateRegistrationRequestRegistrationReference{}
	awsRegistrationReference := &dto.CreateRegistrationRequestRegistrationReference{
		Source: ares.ProviderAWS,
		Token:  testhelper.Randumb(randomdata.Noun()),
	}
	user1 := s.factory.InvitedUser(false)
	org1 := s.factory.Organization(false)
	params1 := &dto.CreateRegistrationRequest{
		Email:                   (*strfmt.Email)(&user1.Email.String),
		OrganizationName:        s.helper.StringPtr(org1.Name),
		OrganizationDisplayName: s.helper.StringPtr(org1.DisplayName),
		RegistrationReference:   registrationReference,
	}

	user2 := s.factory.InvitedUser(true)
	org2 := s.factory.Organization(false)
	params2 := &dto.CreateRegistrationRequest{
		Email:                   (*strfmt.Email)(&user2.Email.String),
		OrganizationName:        s.helper.StringPtr(org2.Name),
		OrganizationDisplayName: s.helper.StringPtr(org2.DisplayName),
		RegistrationReference:   registrationReference,
	}

	user3 := s.factory.InvitedUser(false)
	org3 := s.factory.Organization(true)
	params3 := &dto.CreateRegistrationRequest{
		Email:                   (*strfmt.Email)(&user3.Email.String),
		OrganizationName:        s.helper.StringPtr(org3.Name),
		OrganizationDisplayName: s.helper.StringPtr(org3.DisplayName),
		RegistrationReference:   registrationReference,
	}

	user4 := s.factory.InvitedUser(false)
	org4 := s.factory.Organization(false)
	params4 := &dto.CreateRegistrationRequest{
		Email:                   (*strfmt.Email)(&user4.Email.String),
		OrganizationName:        s.helper.StringPtr(org4.Name),
		OrganizationDisplayName: s.helper.StringPtr(org4.DisplayName),
		RegistrationReference:   awsRegistrationReference,
	}

	user5 := s.factory.InvitedUser(false)
	org5 := s.factory.Organization(false)
	params5 := &dto.CreateRegistrationRequest{
		Email:                   (*strfmt.Email)(&user5.Email.String),
		OrganizationName:        s.helper.StringPtr(org5.Name),
		OrganizationDisplayName: s.helper.StringPtr(org5.DisplayName),
		RegistrationReference:   awsRegistrationReference,
	}

	user5Other := s.factory.InvitedUser(false)
	org5Other := s.factory.Organization(false)
	params5Other := &dto.CreateRegistrationRequest{
		Email:                   (*strfmt.Email)(&user5Other.Email.String),
		OrganizationName:        s.helper.StringPtr(org5Other.Name),
		OrganizationDisplayName: s.helper.StringPtr(org5Other.DisplayName),
		RegistrationReference:   awsRegistrationReference,
	}

	customerID5 := testhelper.Randumb(randomdata.Noun())
	productCode5 := testhelper.Randumb(randomdata.Noun())

	mock := mock_ares.NewMockAWSClient(gomock.NewController(s.T()))
	mock.EXPECT().ResolveCustomer(params5Other.RegistrationReference.Token).Return(customerID5, productCode5, nil)
	s.service.(*auth.Service).AWS = mock
	// when
	err := s.service.CreateRegistration(params5Other)
	xrequire.NoError(s.T(), err)

	cases := map[string]struct {
		MockAWSToken    interface{}
		MockCustomerID  string
		MockProductCode string
		MockError       error
		Params          *dto.CreateRegistrationRequest
		Success         bool
	}{
		"successful registration - default plan": {
			params1.RegistrationReference.Token,
			testhelper.Randumb(randomdata.Noun()),
			testhelper.Randumb(randomdata.Noun()),
			nil,
			params1,
			true,
		},
		"email already used makes registration fail": {
			params2.RegistrationReference.Token,
			testhelper.Randumb(randomdata.Noun()),
			testhelper.Randumb(randomdata.Noun()),
			nil,
			params2,
			false,
		},
		"orgname already used makes registration fail": {
			params3.RegistrationReference.Token,
			testhelper.Randumb(randomdata.Noun()),
			testhelper.Randumb(randomdata.Noun()),
			nil,
			params3,
			false,
		},
		"signup through aws adds billing provider and subscription": {
			params4.RegistrationReference.Token,
			testhelper.Randumb(randomdata.Noun()),
			testhelper.Randumb(randomdata.Noun()),
			nil,
			params4,
			true,
		},
		"fail with duplicate billing info": {
			params5.RegistrationReference.Token,
			customerID5,
			productCode5,
			nil,
			params5,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			mock := mock_ares.NewMockAWSClient(gomock.NewController(t))
			mock.EXPECT().ResolveCustomer(tc.MockAWSToken).Return(tc.MockCustomerID, tc.MockProductCode, tc.MockError)
			s.service.(*auth.Service).AWS = mock
			// when
			err := s.service.CreateRegistration(tc.Params)

			// then
			if tc.Success {
				xrequire.NoError(t, err)

				userFromDB := s.helper.DBGetUserByEmail(tc.Params.Email.String())
				require.Equal(t, null.StringFrom(ares.StatusPendingConfirmation), userFromDB.Status)

				orgFromDB := s.helper.DBGetOrgByName(*tc.Params.OrganizationName)
				require.True(t, time.Now().UTC().After(orgFromDB.R.Subscriptions[0].StartDate))

				if tc.Params.RegistrationReference.Source == "" {
					require.Equal(t, 1, len(orgFromDB.R.Subscriptions))
					plan := s.helper.DBGetSubscriptionPlan(orgFromDB.R.Subscriptions[0].SubscriptionPlanID)
					require.Equal(t, ares.TierFreePlan, plan.Name)
					require.Equal(t, 1, len(orgFromDB.R.OrganizationBillingProviders))
					require.Equal(t, ares.ProviderUnchain, orgFromDB.R.OrganizationBillingProviders[0].ProviderName)
				}
				if tc.Params.RegistrationReference.Source == ares.ProviderAWS {
					require.Equal(t, ares.ProviderAWS, orgFromDB.R.OrganizationBillingProviders[0].ProviderName)
					require.NotNil(t, orgFromDB.R.OrganizationBillingProviders[0].BillingInfo)
				}

				ac := s.helper.DBGetAccountConfirmation(userFromDB.ID)

				require.True(t, ac.ExpirationTime.Before(time.Now().UTC().Add(49*time.Hour)))
				require.True(t, ac.ExpirationTime.After(time.Now().UTC().Add(47*time.Hour)))

				mail := s.helper.GetSignUpCodeFromEmail(ac.Token, tc.Params.Email.String())
				require.NotNil(t, mail)

			} else {
				xrequire.Error(t, err)
			}
		})
	}
}
