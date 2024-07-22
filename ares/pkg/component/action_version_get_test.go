package component_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_GetActionVersion() {
	actionVersion1, action1, org1, _ := s.factory.ActionVersionOrgUser(false, true)
	_, action2, org2, _ := s.factory.ActionVersionOrgUser(false, true)
	actionVersion3, action3 := s.factory.ActionVersionAndAction(false, true)
	actionVersion4, _, org4, _ := s.factory.ActionVersionOrgUser(false, true)
	actionVersion5, action5 := s.factory.ActionVersionAndAction(false, true)

	cases := map[string]struct {
		ActionVersion *orm.ActionVersion
		Action        *orm.Action
		Organization  *orm.Organization
		Success       bool
	}{
		"get action version": {
			actionVersion1,
			action1,
			org1,
			true,
		},
		"version does not exist": {
			s.factory.ActionVersion(false, false),
			action2,
			org2,
			false,
		},
		"organization does not exist": {
			actionVersion3,
			action3,
			s.factory.Organization(false),
			false,
		},
		"action does not exist": {
			actionVersion4,
			s.factory.Action(false, false),
			org4,
			false,
		},
		"organization not creator of action": {
			actionVersion5,
			action5,
			s.factory.Organization(true),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetActionVersion(tc.Organization.Name, tc.Action.Name, tc.ActionVersion.Version)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.ActionVersion.Version, *response.Version)
				require.Equal(t, tc.ActionVersion.ID, *response.ID)
			} else {
				xrequire.Error(t, err)
			}
		})
	}
}

func (s *TestSuite) TestService_GetPublicActionVersion() {
	actionVersion1, action1 := s.factory.ActionVersionAndAction(true, true)
	_, action2 := s.factory.ActionVersionAndAction(true, true)
	actionVersion3, _ := s.factory.ActionVersionAndAction(true, true)

	cases := map[string]struct {
		ActionVersion *orm.ActionVersion
		Action        *orm.Action
		Success       bool
	}{
		"get public action version": {
			actionVersion1,
			action1,
			true,
		},
		"public version does not exist": {
			s.factory.ActionVersion(true, false),
			action2,
			false,
		},
		"public action does not exist": {
			actionVersion3,
			s.factory.Action(true, false),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetPublicActionVersion(tc.Action.Name, tc.ActionVersion.Version)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.ActionVersion.Version, *response.Version)
				require.Equal(t, tc.ActionVersion.ID, *response.ID)
			} else {
				xrequire.Error(t, err)
			}
		})
	}
}
