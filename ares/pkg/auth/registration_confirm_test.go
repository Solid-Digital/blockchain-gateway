package auth_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null"
)

func (s *TestSuite) TestService_ConfirmUserRegistration() {
	password := "some-pwd123"
	user1, ac1 := s.factory.RegisteredUser()
	params1 := &dto.ConfirmRegistrationRequest{
		Token:    ac1.Token,
		Password: password,
		FullName: user1.FullName,
	}

	user2, _ := s.factory.RegisteredUser()
	params2 := &dto.ConfirmRegistrationRequest{
		Token:    "",
		FullName: user2.FullName,
		Password: password,
	}

	user3 := s.factory.InvitedUser(false)
	params3 := &dto.ConfirmRegistrationRequest{
		Token:    "12345",
		Password: password,
		FullName: user3.FullName,
	}

	cases := map[string]struct {
		Params  *dto.ConfirmRegistrationRequest
		Success bool
	}{
		"valid sign up token":      {params1, true},
		"invalid sign up token":    {params2, false},
		"no registration for user": {params3, false},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.ConfirmRegistration(tc.Params)

			if tc.Success {
				xrequire.NoError(t, err)
				userFromAuth, err := s.service.Authenticate(response.Token)
				require.NoError(t, err)
				userFromDB := s.helper.DBGetUser(userFromAuth.ID)

				require.NotNil(t, response.Token)
				require.Equal(t, tc.Params.FullName, userFromDB.FullName)
				xrequire.NoError(t, s.service.CompareHashAndPassword(userFromDB.PasswordHash, password))
				require.Equal(t, null.StringFrom(ares.StatusActive), userFromDB.Status)
			} else {
				xrequire.Error(t, err)
			}
		})
	}

}
