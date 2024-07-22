package component_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_GetAllTriggers() {
	preExistingPublicTriggersCount := len(s.helper.DBPublicTriggers())

	org1, _ := s.factory.OrganizationAndUser(true)
	var included1 = make([]*orm.Trigger, 5)
	for i := 0; i < 5; i++ {
		included1[i] = s.factory.TriggerForOrg(org1, false, true)
	}

	var excluded1 = make([]*orm.Trigger, 5)
	for i := 0; i < 5; i++ {
		excluded1[i] = s.factory.Trigger(false, true)
	}

	org2, _ := s.factory.OrganizationAndUser(true)
	var included2 = make([]*orm.Trigger, 5)
	for i := 0; i < 5; i++ {
		included2[i] = s.factory.TriggerForOrg(org2, false, true)
	}

	cases := map[string]struct {
		Organization        *orm.Organization
		Available           *bool
		Included            []*orm.Trigger
		Excluded            []*orm.Trigger
		PublicTriggerCount  int
		PrivateTriggerCount int
		Success             bool
	}{
		"should return triggers from org1": {
			org1,
			nil,
			included1,
			excluded1,
			preExistingPublicTriggersCount,
			5,
			true,
		},
		"with available true": {
			org1,
			s.helper.BoolPtr(true),
			included1,
			excluded1,
			preExistingPublicTriggersCount,
			5,
			true,
		},
		"with available false": {
			org1,
			s.helper.BoolPtr(false),
			[]*orm.Trigger{},
			[]*orm.Trigger{},
			0, // is this the expected behaviour?
			0, // is this the expected behaviour?
			true,
		},
		"should return triggers from org2": {
			org2,
			nil,
			included2,
			excluded1,
			preExistingPublicTriggersCount,
			5,
			true,
		},
		"organization has no triggers": {
			s.factory.Organization(true),
			nil,
			[]*orm.Trigger{},
			[]*orm.Trigger{},
			preExistingPublicTriggersCount,
			0,
			true,
		},
		"organization does not exist": {
			s.factory.Organization(false),
			nil,
			[]*orm.Trigger{},
			[]*orm.Trigger{},
			preExistingPublicTriggersCount,
			0,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetAllTriggers(tc.Organization.Name, tc.Available)

			if tc.Success {
				xrequire.NoError(t, err)
				publicCount, privateCount := countTriggers(response)
				require.Equal(t, tc.PublicTriggerCount, publicCount)
				require.Equal(t, tc.PrivateTriggerCount, privateCount)
				require.True(t, containsAllTriggers(response, tc.Included))
				require.True(t, containsNoneTriggers(response, tc.Excluded))
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}

func (s *TestSuite) TestService_GetAllPublicTriggers() {
	org1, _ := s.factory.OrganizationAndUser(true)
	s.factory.TriggerForOrg(org1, true, true)

	preExistingPublicTriggersCount := len(s.helper.DBPublicTriggers())

	cases := map[string]struct {
		PublicTriggerCount int
		Success            bool
	}{
		"should return public triggers": {
			preExistingPublicTriggersCount,
			true,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetAllPublicTriggers()

			if tc.Success {
				xrequire.NoError(t, err)
				publicCount, _ := countTriggers(response)
				require.Equal(t, tc.PublicTriggerCount, publicCount)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}

func countTriggers(responses []*dto.GetComponentResponse) (publicCount int, privateCount int) {
	for _, res := range responses {
		if *res.Public {
			publicCount++
		} else {
			privateCount++
		}
	}

	return publicCount, privateCount
}

func containsTriggerID(lst []*dto.GetComponentResponse, ID int64) bool {
	for _, e := range lst {
		if *e.ID == ID {
			return true
		}
	}

	return false
}

func containsAllTriggers(response []*dto.GetComponentResponse, triggers []*orm.Trigger) bool {
	for _, trigger := range triggers {
		containsID := containsTriggerID(response, trigger.ID)

		if !containsID {
			return false
		}
	}

	return true
}

func containsNoneTriggers(response []*dto.GetComponentResponse, triggers []*orm.Trigger) bool {
	for _, trigger := range triggers {
		containsID := containsTriggerID(response, trigger.ID)

		if containsID {
			return false
		}
	}

	return true
}
