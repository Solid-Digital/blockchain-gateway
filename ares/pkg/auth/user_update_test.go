package auth_test

import (
	"testing"

	"github.com/go-openapi/strfmt"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/dto"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_UpdateCurrentUser() {
	cases := map[string]struct {
		Params  *dto.UpdateCurrentUserRequest
		User    *dto.User
		Success bool
	}{
		"change full name": {
			&dto.UpdateCurrentUserRequest{
				FullName: s.factory.User(false).FullName,
			},
			s.factory.DTOUser(true),
			true},
		"change email": {
			&dto.UpdateCurrentUserRequest{
				Email: strfmt.Email(s.factory.User(false).Email.String),
			},
			s.factory.DTOUser(true),
			true},
		"change email and full name": {
			&dto.UpdateCurrentUserRequest{
				Email:    strfmt.Email(s.factory.User(false).Email.String),
				FullName: s.factory.User(false).FullName,
			},
			s.factory.DTOUser(true),
			true},
		"user does not exist": {
			&dto.UpdateCurrentUserRequest{
				FullName: s.factory.User(false).FullName,
			},
			s.factory.DTOUser(false),
			false},
		"test duplicate uppercase email": {
			&dto.UpdateCurrentUserRequest{
				Email: strfmt.Email(s.factory.User(true).Email.String),
			},
			s.factory.DTOUser(true),
			false},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, appErr := s.service.UpdateCurrentUser(tc.Params, tc.User)

			if tc.Success {
				xrequire.NoError(t, appErr)
				if tc.Params.Email != "" {
					require.Equal(t, string(tc.Params.Email), response.Email)
					require.Equal(t, string(tc.Params.Email), s.helper.DBGetUser(tc.User.ID).Email.String)
				}
				if tc.Params.FullName != "" {
					require.Equal(t, tc.Params.FullName, response.FullName)
					require.Equal(t, tc.Params.FullName, s.helper.DBGetUser(tc.User.ID).FullName)
				}
			} else {
				xrequire.Error(t, appErr)
				require.Nil(t, response)
			}
		})
	}
}
