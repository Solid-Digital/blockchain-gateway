package component_test

import (
	"testing"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_GetTrigger() {
	trigger1, org1, _ := s.factory.TriggerOrgUser(false, true)

	trigger2, org2, _ := s.factory.TriggerOrgUser(true, true)
	_ = s.factory.TriggerVersionForOrgAndTrigger(false, true, org2, trigger2, "v0.0.1")
	_ = s.factory.TriggerVersionForOrgAndTrigger(true, true, org2, trigger2, "v0.0.100")
	_ = s.factory.TriggerVersionForOrgAndTrigger(false, true, org2, trigger2, "v0.0.11")
	trigger3, org3, _ := s.factory.TriggerOrgUser(true, true)
	_ = s.factory.TriggerVersionForOrgAndTrigger(true, true, org3, trigger3, "v0.0.1")

	cases := map[string]struct {
		Trigger          *orm.Trigger
		Organization     *orm.Organization
		ExpectedVersions []string
		Success          bool
	}{
		"get trigger": {
			trigger1,
			org1,
			nil,
			true,
		},
		"get trigger returns sorted versions": {
			trigger2,
			org2,
			[]string{"v0.0.100", "v0.0.11", "v0.0.1"},
			true,
		},
		"get trigger doesn't panic if there are private trigger versions in another organization": {
			trigger2,
			org3,
			[]string{"v0.0.100"},
			true,
		},
		"trigger does not exist": {
			s.factory.Trigger(false, false),
			s.factory.Organization(true),
			nil,
			false,
		},
		"organization does not exist": {
			s.factory.Trigger(false, true),
			s.factory.Organization(false),
			nil,
			false,
		},
		"trigger not related to organization": {
			s.factory.Trigger(false, true),
			s.factory.Organization(true),
			nil,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetTrigger(tc.Organization.Name, tc.Trigger.Name)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.Trigger.Name, *response.Name)
				require.Equal(t, tc.ExpectedVersions, response.Versions)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}

func (s *TestSuite) TestService_GetPublicTrigger() {
	trigger1, org1, _ := s.factory.TriggerOrgUser(true, true)
	_ = s.factory.TriggerVersionForOrgAndTrigger(true, true, org1, trigger1, "v0.0.1")

	cases := map[string]struct {
		Trigger          *orm.Trigger
		Organization     *orm.Organization
		ExpectedVersions []string
		Success          bool
	}{
		"get public trigger": {
			trigger1,
			org1,
			[]string{"v0.0.1"},
			true,
		},
		"public trigger does not exist": {
			s.factory.Trigger(false, false),
			s.factory.Organization(true),
			nil,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetPublicTrigger(tc.Trigger.Name)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.Trigger.Name, *response.Name)
				require.Equal(t, tc.ExpectedVersions, response.Versions)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}
