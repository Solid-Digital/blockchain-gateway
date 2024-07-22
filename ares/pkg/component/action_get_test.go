package component_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_GetAction() {
	action1, org1, _ := s.factory.ActionOrgUser(false, true)

	action2, org2, _ := s.factory.ActionOrgUser(true, true)
	_ = s.factory.ActionVersionForOrgAndAction(false, true, org2, action2, "v0.0.1")
	_ = s.factory.ActionVersionForOrgAndAction(true, true, org2, action2, "v0.0.100")
	_ = s.factory.ActionVersionForOrgAndAction(false, true, org2, action2, "v0.0.11")
	action3, org3, _ := s.factory.ActionOrgUser(true, true)
	_ = s.factory.ActionVersionForOrgAndAction(true, true, org3, action3, "v0.0.1")

	cases := map[string]struct {
		Action           *orm.Action
		Organization     *orm.Organization
		ExpectedVersions []string
		Success          bool
	}{
		"get action": {
			action1,
			org1,
			nil,
			true,
		},
		"get action returns sorted versions": {
			action2,
			org2,
			[]string{"v0.0.100", "v0.0.11", "v0.0.1"},
			true,
		},
		"get action doesn't panic if there are private action versions in another organization": {
			action2,
			org3,
			[]string{"v0.0.100"},
			true,
		},
		"action does not exist": {
			s.factory.Action(false, false),
			s.factory.Organization(true),
			nil,
			false,
		},
		"organization does not exist": {
			s.factory.Action(false, true),
			s.factory.Organization(false),
			nil,
			false,
		},
		"action not related to organization": {
			s.factory.Action(false, true),
			s.factory.Organization(true),
			nil,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetAction(tc.Organization.Name, tc.Action.Name)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.Action.Name, *response.Name)
				require.Equal(t, tc.ExpectedVersions, response.Versions)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}

func (s *TestSuite) TestService_GetPublicAction() {
	action1, org1, _ := s.factory.ActionOrgUser(true, true)
	_ = s.factory.ActionVersionForOrgAndAction(true, true, org1, action1, "v0.0.1")

	cases := map[string]struct {
		Action           *orm.Action
		Organization     *orm.Organization
		ExpectedVersions []string
		Success          bool
	}{
		"get action": {
			action1,
			org1,
			[]string{"v0.0.1"},
			true,
		},
		"action does not exist": {
			s.factory.Action(false, false),
			s.factory.Organization(true),
			nil,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetPublicAction(tc.Action.Name)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.Action.Name, *response.Name)
				require.Equal(t, tc.ExpectedVersions, response.Versions)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}
