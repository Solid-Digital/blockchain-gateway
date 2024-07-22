package auth_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
)

func (s *TestSuite) TestService_GetCurrentUser() {
	_, user := s.factory.OrganizationAndUser(true)

	cases := map[string]struct {
		User    *dto.User
		Success bool
	}{
		"user exists": {
			s.factory.ORMToDTOUser(user),
			true},
		"user does not exist": {
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, appErr := s.service.GetCurrentUser(tc.User)

			if tc.Success {
				xrequire.NoError(t, appErr)
				require.Equal(t, string(tc.User.Email), response.Email)
			} else {
				xrequire.Error(t, appErr)
				require.Nil(t, response)
			}
		})
	}
}
