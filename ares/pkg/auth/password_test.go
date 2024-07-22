package auth_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
)

func (s *TestSuite) TestService_HashPassword() {
	password := "qwerty85"
	passwordHash, appErr := s.service.HashPassword(password)

	xrequire.NoError(s.T(), appErr)
	require.True(s.T(), s.helper.CorrectPassword(password, passwordHash))
}

func (s *TestSuite) TestService_CorrectHashAndPassword() {
	cases := map[string]struct {
		Hash     string
		Password string
		Success  bool
	}{
		"valid password": {
			"$2a$10$nhl0yYEeN3QEWH/vXf2ZDOF27wyRW1KiMK5HR9FG22rvh8IByDPhu",
			"qwerty85",
			true},
		"invalid password": {
			"$2a$10$nhl0yYEeN3QEWH/vXf2ZDOF27wyRW1KiMK5HR9FG22rvh8IByDPhu",
			"invalid-123",
			false},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			err := s.service.CompareHashAndPassword(tc.Hash, tc.Password)

			if tc.Success {
				xrequire.NoError(t, err)
			} else {
				xrequire.Error(t, err)
			}
		})
	}
}
