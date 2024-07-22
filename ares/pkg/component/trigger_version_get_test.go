package component_test

import (
	"testing"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_GetTriggerVersion() {
	triggerVersion1, trigger1, org1, _ := s.factory.TriggerVersionOrgUser(false, true)
	_, trigger2, org2, _ := s.factory.TriggerVersionOrgUser(false, true)
	triggerVersion3, trigger3 := s.factory.TriggerVersionAndTrigger(false, true)
	triggerVersion4, _, org4, _ := s.factory.TriggerVersionOrgUser(false, true)
	triggerVersion5, trigger5 := s.factory.TriggerVersionAndTrigger(false, true)

	cases := map[string]struct {
		TriggerVersion *orm.TriggerVersion
		Trigger        *orm.Trigger
		Organization   *orm.Organization
		Success        bool
	}{
		"get trigger version": {
			triggerVersion1,
			trigger1,
			org1,
			true,
		},
		"version does not exist": {
			s.factory.TriggerVersion(false, false),
			trigger2,
			org2,
			false,
		},
		"organization does not exist": {
			triggerVersion3,
			trigger3,
			s.factory.Organization(false),
			false,
		},
		"trigger does not exist": {
			triggerVersion4,
			s.factory.Trigger(false, false),
			org4,
			false,
		},
		"organization not creator of trigger": {
			triggerVersion5,
			trigger5,
			s.factory.Organization(true),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetTriggerVersion(tc.Organization.Name, tc.Trigger.Name, tc.TriggerVersion.Version)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.TriggerVersion.Version, *response.Version)
				require.Equal(t, tc.TriggerVersion.ID, *response.ID)
			} else {
				xrequire.Error(t, err)
			}
		})
	}
}

func (s *TestSuite) TestService_GetPublicTriggerVersion() {
	triggerVersion1, trigger1 := s.factory.TriggerVersionAndTrigger(true, true)
	_, trigger2 := s.factory.TriggerVersionAndTrigger(true, true)
	triggerVersion3, _ := s.factory.TriggerVersionAndTrigger(true, true)

	cases := map[string]struct {
		TriggerVersion *orm.TriggerVersion
		Trigger        *orm.Trigger
		Success        bool
	}{
		"get trigger version": {
			triggerVersion1,
			trigger1,
			true,
		},
		"version does not exist": {
			s.factory.TriggerVersion(true, false),
			trigger2,
			false,
		},
		"trigger does not exist": {
			triggerVersion3,
			s.factory.Trigger(true, false),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetPublicTriggerVersion(tc.Trigger.Name, tc.TriggerVersion.Version)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.TriggerVersion.Version, *response.Version)
				require.Equal(t, tc.TriggerVersion.ID, *response.ID)
			} else {
				xrequire.Error(t, err)
			}
		})
	}
}
