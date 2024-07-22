package component_test

import (
	"testing"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_GetBaseVersion() {
	baseVersion1, base1, org1, _ := s.factory.BaseVersionOrgUser(false, true)
	_, base2, org2, _ := s.factory.BaseVersionOrgUser(false, true)
	baseVersion3, base3 := s.factory.BaseVersionAndBase(false, true)
	baseVersion4, _, org4, _ := s.factory.BaseVersionOrgUser(false, true)
	baseVersion5, base5 := s.factory.BaseVersionAndBase(false, true)

	cases := map[string]struct {
		BaseVersion  *orm.BaseVersion
		Base         *orm.Base
		Organization *orm.Organization
		Success      bool
	}{
		"get base version": {
			baseVersion1,
			base1,
			org1,
			true,
		},
		"version does not exist": {
			s.factory.BaseVersion(false, false),
			base2,
			org2,
			false,
		},
		"organization does not exist": {
			baseVersion3,
			base3,
			s.factory.Organization(false),
			false,
		},
		"base does not exist": {
			baseVersion4,
			s.factory.Base(false, false),
			org4,
			false,
		},
		"organization not creator of base": {
			baseVersion5,
			base5,
			s.factory.Organization(true),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetBaseVersion(tc.Organization.Name, tc.Base.Name, tc.BaseVersion.Version)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.BaseVersion.Version, *response.Version)
				require.Equal(t, tc.BaseVersion.ID, *response.ID)
			} else {
				xrequire.Error(t, err)
			}
		})
	}
}
