package auth_test

import (
	"testing"

	"github.com/volatiletech/null"

	"bitbucket.org/unchain/ares/pkg/auth"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/dto"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_DeleteCurrentUser() {
	cases := map[string]struct {
		User    *dto.User
		Success bool
	}{
		"user exists": {
			s.factory.DTOUser(true),
			true},
		"user does not exist": {
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			appErr := s.service.DeleteCurrentUser(tc.User)

			if tc.Success {
				xrequire.NoError(t, appErr)
				require.True(t, s.helper.DBUserExists(tc.User.ID))
				require.Equal(t, s.helper.DBGetUser(tc.User.ID).Email, null.StringFromPtr(nil))
				require.Equal(t, s.helper.DBGetUser(tc.User.ID).FullName, auth.ArchivedFullName)
			} else {
				xrequire.Error(t, appErr)
			}
		})
	}
}
