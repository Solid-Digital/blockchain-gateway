package auth_test

import (
	"testing"

	"github.com/go-openapi/strfmt"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"github.com/stretchr/testify/require"

	"bitbucket.org/unchain/ares/gen/dto"
)

const (
	recoveryFmt = "recovery:%s"
)

func (s *TestSuite) TestService_ResetPassword() {
	cases := map[string]struct {
		Params  *dto.ResetPasswordRequest
		Success bool
	}{
		"user exists": {
			&dto.ResetPasswordRequest{
				Email: strfmt.Email(s.factory.User(true).Email.String)},
			true},
		"user does not exist": {
			&dto.ResetPasswordRequest{
				Email: strfmt.Email(s.factory.User(false).Email.String)},
			false},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.ResetPassword(tc.Params)

			if tc.Success {
				xrequire.NoError(t, err)
				require.NotEmpty(t, s.helper.GetRecoveryCodeFromEmail(response.RequestID, tc.Params.Email.String()))
			} else {
				xrequire.Error(t, err)
				require.Empty(t, s.helper.GetMailbox(tc.Params.Email.String()))
			}
		})
	}
}
