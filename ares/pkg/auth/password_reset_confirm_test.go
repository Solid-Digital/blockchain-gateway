package auth_test

import (
	"testing"

	"bitbucket.org/unchain/ares/gen/orm"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"github.com/stretchr/testify/require"

	"bitbucket.org/unchain/ares/gen/dto"
)

func (s *TestSuite) TestService_ConfirmResetPassword() {
	user1 := s.factory.User(true)
	user2 := s.factory.User(true)

	cases := map[string]struct {
		Params  *dto.ConfirmResetPasswordRequest
		User    *orm.User
		Success bool
	}{
		"valid code": {
			&dto.ConfirmResetPasswordRequest{
				Password:     newPassword,
				RecoveryCode: s.factory.RecoveryCode(user1.Email.String),
			},
			user1,
			true},
		"invalid code": {
			&dto.ConfirmResetPasswordRequest{
				Password:     newPassword,
				RecoveryCode: "12345",
			},
			user2,
			false},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			appErr := s.service.ConfirmResetPassword(tc.Params)

			if tc.Success {
				xrequire.NoError(t, appErr)
				require.True(t, s.helper.CorrectUserPassword(tc.User.ID, newPassword))
			} else {
				xrequire.Error(t, appErr)
				require.False(t, s.helper.CorrectUserPassword(tc.User.ID, newPassword))
			}
		})
	}
}
