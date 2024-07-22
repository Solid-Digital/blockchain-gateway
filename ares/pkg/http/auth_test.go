package http_test

import (
	"math/rand"
	stdhttp "net/http"
	"testing"

	"bitbucket.org/unchain/ares/gen/api/operations/auth"

	"bitbucket.org/unchain/ares/gen/dto"
	"github.com/stretchr/testify/require"
	"github.com/unchainio/pkg/errors"

	"github.com/golang/mock/gomock"

	mock_ares "bitbucket.org/unchain/ares/gen/mocks"

	"bitbucket.org/unchain/ares/pkg/factory"
	"bitbucket.org/unchain/ares/pkg/http"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	factory       *factory.Factory
	handler       *http.AuthHandler
	mockedService *mock_ares.MockAuthService
}

func (s *AuthTestSuite) SetupSuite() {
	s.T().Skip()
	s.factory = factory.NewFactory(&s.Suite)
}

// This runs before each test
// A new mock needs to be created for each test, otherwise it will fail
func (s *AuthTestSuite) SetupTest() {
	service := mock_ares.NewMockAuthService(gomock.NewController(s.T()))
	s.mockedService = service
	s.handler = http.NewAuthHandler(service, s.factory.Logger())
}

func (s *AuthTestSuite) TestAuthHandler_ChangeCurrentPassword() {
	cases := map[string]struct {
		ServiceError error
		Params       auth.ChangeCurrentPasswordParams
		Principal    *dto.User
		Success      bool
	}{
		"service returns no error": {
			nil,
			auth.NewChangeCurrentPasswordParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			errors.New("failed"),
			auth.NewChangeCurrentPasswordParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().ChangeCurrentPassword(gomock.Any(), gomock.Any()).Return(tc.ServiceError)
			response := s.handler.ChangeCurrentPassword(tc.Params, tc.Principal)

			if tc.Success {
				_, ok := response.(*auth.ChangeCurrentPasswordCreated)
				require.True(t, ok)
			} else {
				result, ok := response.(*auth.ChangeCurrentPasswordUnauthorized)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *AuthTestSuite) TestAuthHandler_ConfirmResetPassword() {
	cases := map[string]struct {
		ServiceError error
		Params       auth.ConfirmResetPasswordParams
		Success      bool
	}{
		"service returns no error": {
			nil,
			auth.NewConfirmResetPasswordParams(),
			true},
		"service returns an error": {
			errors.New("failed"),
			auth.NewConfirmResetPasswordParams(),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().ConfirmResetPassword(gomock.Any()).Return(tc.ServiceError)
			response := s.handler.ConfirmResetPassword(tc.Params)

			if tc.Success {
				_, ok := response.(*auth.ConfirmResetPasswordOK)
				require.True(t, ok)
			} else {
				result, ok := response.(*auth.ConfirmResetPasswordUnauthorized)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *AuthTestSuite) TestAuthHandler_CreateRegistration() {
	cases := map[string]struct {
		ServiceReturn *auth.CreateRegistrationOK
		ServiceError  error
		Params        auth.CreateRegistrationParams
		Success       bool
	}{
		"service returns no error": {
			&auth.CreateRegistrationOK{},
			nil,
			auth.NewCreateRegistrationParams(),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			auth.NewCreateRegistrationParams(),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().CreateRegistration(gomock.Any()).Return(tc.ServiceError)
			response := s.handler.CreateRegistration(tc.Params)

			if tc.Success {
				result, ok := response.(*auth.CreateRegistrationOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result)
			} else {
				result, ok := response.(*auth.CreateRegistrationInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *AuthTestSuite) TestAuthHandler_ConfirmRegistration() {
	cases := map[string]struct {
		ServiceReturn *dto.LoginResponse
		ServiceError  error
		Params        auth.ConfirmRegistrationParams
		Success       bool
	}{
		"service returns no error": {
			&dto.LoginResponse{UserID: rand.Int63()},
			nil,
			auth.NewConfirmRegistrationParams(),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			auth.NewConfirmRegistrationParams(),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().ConfirmRegistration(gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.ConfirmRegistration(tc.Params)

			if tc.Success {
				result, ok := response.(*auth.ConfirmRegistrationOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*auth.ConfirmRegistrationInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *AuthTestSuite) TestAuthHandler_DeleteCurrentUser() {
	cases := map[string]struct {
		ServiceError error
		Params       auth.DeleteCurrentUserParams
		Principal    *dto.User
		Success      bool
	}{
		"service returns no error": {
			nil,
			auth.NewDeleteCurrentUserParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			errors.New("failed"),
			auth.NewDeleteCurrentUserParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().DeleteCurrentUser(gomock.Any()).Return(tc.ServiceError)
			response := s.handler.DeleteCurrentUser(tc.Params, tc.Principal)

			if tc.Success {
				_, ok := response.(*auth.DeleteCurrentUserNoContent)
				require.True(t, ok)
			} else {
				result, ok := response.(*auth.DeleteCurrentUserInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *AuthTestSuite) TestAuthHandler_GetCurrentUser() {
	cases := map[string]struct {
		ServiceReturn *dto.GetCurrentUserResponse
		ServiceError  error
		Params        auth.GetCurrentUserParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetCurrentUserResponse{ID: rand.Int63()},
			nil,
			auth.NewGetCurrentUserParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			auth.NewGetCurrentUserParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetCurrentUser(gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetCurrentUser(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*auth.GetCurrentUserOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*auth.GetCurrentUserInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *AuthTestSuite) TestAuthHandler_Login() {
	cases := map[string]struct {
		ServiceReturn *dto.LoginResponse
		ServiceError  error
		Params        auth.LoginParams
		Success       bool
	}{
		"service returns no error": {
			&dto.LoginResponse{UserID: rand.Int63()},
			nil,
			auth.LoginParams{
				HTTPRequest: &stdhttp.Request{RemoteAddr: "example.com:1234"},
			},
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			auth.LoginParams{
				HTTPRequest: &stdhttp.Request{RemoteAddr: "example.com:1234"},
			},
			false},
		"port is missing": {
			&dto.LoginResponse{UserID: rand.Int63()},
			nil,
			auth.LoginParams{
				HTTPRequest: &stdhttp.Request{RemoteAddr: "example.com"},
			},
			true},
		"invalid host": {
			&dto.LoginResponse{UserID: rand.Int63()},
			nil,
			auth.LoginParams{
				HTTPRequest: &stdhttp.Request{RemoteAddr: "example.com::1234"},
			},
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().Login(gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.Login(tc.Params)

			if tc.Success {
				result, ok := response.(*auth.LoginOK)
				require.True(t, ok)
				s.Require().Equal(tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*auth.LoginUnauthorized)
				s.Require().True(ok)

				// I know... this is not great
				if tc.ServiceError != nil {
					require.Contains(t, result.Payload, "failed")
				} else {
					require.NotContains(t, result.Payload, "failed")
				}
			}
		})
	}
}

func (s *AuthTestSuite) TestAuthHandler_Logout() {
	cases := map[string]struct {
		ServiceError error
		Params       auth.LogoutParams
		Principal    *dto.User
		Success      bool
	}{
		"service returns no error": {
			nil,
			auth.NewLogoutParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			errors.New("failed"),
			auth.NewLogoutParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().Logout(gomock.Any()).Return(tc.ServiceError)
			response := s.handler.Logout(tc.Params, tc.Principal)

			if tc.Success {
				_, ok := response.(*auth.LogoutOK)
				require.True(t, ok)
			} else {
				result, ok := response.(*auth.LogoutInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *AuthTestSuite) TestAuthHandler_ResetPassword() {
	cases := map[string]struct {
		ServiceReturn *dto.ResetPasswordResponse
		ServiceError  error
		Params        auth.ResetPasswordParams
		Success       bool
	}{
		"service returns no error": {
			&dto.ResetPasswordResponse{RequestID: string(rand.Int63())},
			nil,
			auth.ResetPasswordParams{
				HTTPRequest: &stdhttp.Request{Header: stdhttp.Header{}},
			},
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			auth.ResetPasswordParams{
				HTTPRequest: &stdhttp.Request{Header: stdhttp.Header{}},
			},
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().ResetPassword(gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.ResetPassword(tc.Params)

			if tc.Success {
				result, ok := response.(*auth.ResetPasswordOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*auth.ResetPasswordUnauthorized)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *AuthTestSuite) TestAuthHandler_UpdateCurrentUser() {
	cases := map[string]struct {
		ServiceReturn *dto.GetCurrentUserResponse
		ServiceError  error
		Params        auth.UpdateCurrentUserParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetCurrentUserResponse{ID: rand.Int63()},
			nil,
			auth.NewUpdateCurrentUserParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			auth.NewUpdateCurrentUserParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().UpdateCurrentUser(gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.UpdateCurrentUser(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*auth.UpdateCurrentUserOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*auth.UpdateCurrentUserInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

// This will make sure the test suite will run
// Don't put any logic for setting up the tests in here, use the hooks from the test suite for that
func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
