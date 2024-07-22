package casbin_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_MakeSuperAdmin() {
	cases := map[string]struct {
		User    *orm.User
		Success bool
	}{
		"valid email":        {s.factory.User(true), true},
		"non existing email": {s.factory.User(false), false},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			err := s.enforcer.MakeSuperAdmin(tc.User.ID)

			spew.Dump(err)

			if tc.Success {
				require.NoError(t, err)
				require.True(t, s.enforcer.GetGlobalRolesForUser(tc.User.ID)[ares.RoleSuperAdmin.String()])

			} else {
				require.NoError(t, err) // is this expected behaviour?

				// TODO: refactor
				//require.False(t, s.getGlobalRoles(tc.User.Email)[ares.SuperAdmin.String()])
			}
		})
	}
}
