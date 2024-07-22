package auth_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
)

var (
	currentPassword = "qwerty85"
	newPassword     = "foobar2000"
)

func (s *TestSuite) TestService_ChangeCurrentPassword() {
	cases := map[string]struct {
		Params  *dto.ChangeCurrentPasswordRequest
		User    *dto.User
		Success bool
	}{
		"password is correct": {
			&dto.ChangeCurrentPasswordRequest{
				CurrentPassword: currentPassword,
				NewPassword:     newPassword,
			},
			s.factory.DTOUser(true),
			true},
		"password not correct": {
			&dto.ChangeCurrentPasswordRequest{
				CurrentPassword: "invalid",
				NewPassword:     newPassword,
			},
			s.factory.DTOUser(true),
			false},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			err := s.service.ChangeCurrentPassword(tc.Params, tc.User)

			if tc.Success {
				xrequire.NoError(t, err)
				require.True(t, s.helper.CorrectUserPassword(tc.User.ID, newPassword))
			} else {
				xrequire.Error(t, err)
				require.False(t, s.helper.CorrectUserPassword(tc.User.ID, newPassword))
			}
		})
	}
}
