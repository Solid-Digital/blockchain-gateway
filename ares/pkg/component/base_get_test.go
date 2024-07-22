package component_test

import (
	"testing"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_GetBase() {
	base1, org1, _ := s.factory.BaseOrgUser(false, true)

	base2, org2, _ := s.factory.BaseOrgUser(true, true)

	_ = s.factory.BaseVersionForOrgAndBase(false, true, org2, base2, "v0.0.1")
	_ = s.factory.BaseVersionForOrgAndBase(true, true, org2, base2, "v0.0.100")
	_ = s.factory.BaseVersionForOrgAndBase(false, true, org2, base2, "v0.0.11")
	base3, org3, _ := s.factory.BaseOrgUser(true, true)
	_ = s.factory.BaseVersionForOrgAndBase(true, true, org3, base3, "v0.0.1")

	cases := map[string]struct {
		Base             *orm.Base
		Organization     *orm.Organization
		ExpectedVersions []string
		Success          bool
	}{
		"get base": {
			base1,
			org1,
			nil,
			true,
		},
		"get base returns sorted versions": {
			base2,
			org2,
			[]string{"v0.0.100", "v0.0.11", "v0.0.1"},
			true,
		},
		"get base doesn't panic if there are private base versions in another organization": {
			base2,
			org3,
			[]string{"v0.0.100"},
			true,
		},
		"base does not exist": {
			s.factory.Base(false, false),
			s.factory.Organization(true),
			nil,
			false,
		},
		"organization does not exist": {
			s.factory.Base(false, true),
			s.factory.Organization(false),
			nil,
			false,
		},
		"base not related to organization": {
			s.factory.Base(false, true),
			s.factory.Organization(true),
			nil,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetBase(tc.Organization.Name, tc.Base.Name)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.Base.Name, *response.Name)
				require.Equal(t, tc.ExpectedVersions, response.Versions)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}
