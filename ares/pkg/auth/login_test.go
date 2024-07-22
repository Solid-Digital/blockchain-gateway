package auth_test

import (
	"strings"
	"testing"
	"time"

	"bitbucket.org/unchain/ares/pkg/auth"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"

	"bitbucket.org/unchain/ares/gen/dto"
	"github.com/Pallinder/go-randomdata"
)

func (s *TestSuite) TestService_Login() {
	cases := map[string]struct {
		IP      string
		Params  *dto.LoginRequest
		Success bool
	}{
		"email and password correct": {
			randomdata.IpV4Address(),
			&dto.LoginRequest{
				Email:    (*strfmt.Email)(s.helper.StringPtr(s.factory.RegisteredUser2().Email.String)),
				Password: s.helper.StringPtr("qwerty85"),
			},
			true},
		"email not case sensitive": {
			randomdata.IpV4Address(),
			&dto.LoginRequest{
				Email:    (*strfmt.Email)(s.helper.StringPtr(strings.ToUpper(s.factory.RegisteredUser2().Email.String))),
				Password: s.helper.StringPtr("qwerty85"),
			},
			true},
		"email not correct": {
			randomdata.IpV4Address(),
			&dto.LoginRequest{
				Email:    (*strfmt.Email)(s.helper.StringPtr("test@example.com")),
				Password: s.helper.StringPtr("qwerty85"),
			},
			false},
		"password not correct": {
			randomdata.IpV4Address(),
			&dto.LoginRequest{
				Email:    (*strfmt.Email)(s.helper.StringPtr(strings.ToUpper(s.factory.RegisteredUser2().Email.String))),
				Password: s.helper.StringPtr("incorrect"),
			},
			false},
		"archived user cannot authenticate": {
			randomdata.IpV4Address(),
			&dto.LoginRequest{
				Email:    (*strfmt.Email)(s.helper.StringPtr(strings.ToUpper(s.factory.ArchivedUser().Email.String))),
				Password: s.helper.StringPtr("qwerty85"),
			},
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, appErr := s.service.Login(tc.IP, tc.Params)

			if tc.Success {
				xrequire.NoError(t, appErr)
				require.Equal(t, s.helper.DBGetUserByEmail(strings.ToLower(tc.Params.Email.String())).ID, response.UserID)
			} else {
				xrequire.Error(t, appErr)
				require.Nil(t, response)
			}
		})
	}
}

// Not a great way of testing, code may require refactoring
func (s *TestSuite) TestService_LoginAttempts() {
	ip := randomdata.IpV4Address()
	user := s.factory.User(true)
	params := &dto.LoginRequest{
		Email:    (*strfmt.Email)(&user.Email.String),
		Password: s.helper.StringPtr("invalid"),
	}

	failOnAttempt := 6

	for i := 1; i <= 10; i++ {
		_, appErr := s.service.Login(ip, params)

		xrequire.Error(s.T(), appErr)

		if i >= failOnAttempt {
			require.Contains(s.T(), appErr.Error(), "blocked")
		} else {
			require.NotContains(s.T(), appErr.Error(), "blocked")
		}
	}
}

func (s *TestSuite) TestService_TokenExpiry() {
	user := s.factory.User(true)

	token, appErr := s.ares.AuthService.(*auth.Service).GenerateToken(user.ID, user.Email.String, 1*time.Second)
	xrequire.NoError(s.T(), appErr)

	time.Sleep(2 * time.Second)

	_, err := s.ares.AuthService.Authenticate(token)
	require.Error(s.T(), err)
}
