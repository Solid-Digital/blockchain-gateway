package auth_test

import (
	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_Logout() {
	token := &dto.Token{
		Expiration: 10,
		Raw:        "foobar",
	}
	appErr := s.service.Logout(token)

	xrequire.NoError(s.T(), appErr)
	require.True(s.T(), s.ares.Redis.IsTokenInBlacklist("foobar"))
}
