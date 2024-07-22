package component_test

import (
	"fmt"
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_GetAllActions() {
	preExistingPublicActionsCount := len(s.helper.DBPublicActions())

	org1, _ := s.factory.OrganizationAndUser(true)
	var included1 = make([]*orm.Action, 5)
	for i := 0; i < 5; i++ {
		included1[i] = s.factory.ActionForOrg(org1, false, true)
	}

	var excluded1 = make([]*orm.Action, 5)
	for i := 0; i < 5; i++ {
		excluded1[i] = s.factory.Action(false, true)
	}

	org2, _ := s.factory.OrganizationAndUser(true)
	var included2 = make([]*orm.Action, 5)
	for i := 0; i < 5; i++ {
		included2[i] = s.factory.ActionForOrg(org2, false, true)
	}

	cases := map[string]struct {
		Organization        *orm.Organization
		Available           *bool
		Included            []*orm.Action
		Excluded            []*orm.Action
		PublicActionsCount  int
		PrivateActionsCount int
		Success             bool
	}{
		"should return actions from org1": {
			org1,
			nil,
			included1,
			excluded1,
			preExistingPublicActionsCount,
			5,
			true,
		},
		"with available true": {
			org1,
			s.helper.BoolPtr(true),
			included1,
			excluded1,
			preExistingPublicActionsCount,
			5,
			true,
		},
		"with available false": {
			org1,
			s.helper.BoolPtr(false),
			[]*orm.Action{},
			[]*orm.Action{},
			0, // is this the expected behaviour?
			0, // is this the expected behaviour?
			true,
		},
		"should return actions from org2": {
			org2,
			nil,
			included2,
			excluded1,
			preExistingPublicActionsCount,
			5,
			true,
		},
		"organization has no actions": {
			s.factory.Organization(true),
			nil,
			[]*orm.Action{},
			[]*orm.Action{},
			preExistingPublicActionsCount,
			0,
			true,
		},
		"organization does not exist": {
			s.factory.Organization(false),
			nil,
			[]*orm.Action{},
			[]*orm.Action{},
			preExistingPublicActionsCount,
			0,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetAllActions(tc.Organization.Name, tc.Available)

			if tc.Success {
				xrequire.NoError(t, err)
				publicCount, privateCount := countActions(response)
				require.Equal(t, tc.PublicActionsCount, publicCount)
				require.Equal(t, tc.PrivateActionsCount, privateCount)
				require.True(t, containsAllActions(response, tc.Included))
				require.True(t, containsNoneActions(response, tc.Excluded))
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}

func (s *TestSuite) TestService_GetAllPublicActions() {
	org1, _ := s.factory.OrganizationAndUser(true)
	s.factory.ActionForOrg(org1, true, true)

	preExistingPublicActionsCount := len(s.helper.DBPublicActions())
	fmt.Println("Pre existing public actions: ", preExistingPublicActionsCount)

	cases := map[string]struct {
		PublicActionsCount int
		Success            bool
	}{
		"should return all public actions": {
			preExistingPublicActionsCount,
			true,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetAllPublicActions()

			if tc.Success {
				xrequire.NoError(t, err)
				publicCount, _ := countActions(response)
				require.Equal(t, tc.PublicActionsCount, publicCount)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}

func countActions(responses []*dto.GetComponentResponse) (publicCount int, privateCount int) {
	for _, res := range responses {
		if *res.Public {
			publicCount++
		} else {
			privateCount++
		}
	}

	return publicCount, privateCount
}

func containsActionID(lst []*dto.GetComponentResponse, ID int64) bool {
	for _, e := range lst {
		if *e.ID == ID {
			return true
		}
	}

	return false
}

func containsAllActions(response []*dto.GetComponentResponse, actions []*orm.Action) bool {
	for _, action := range actions {
		containsID := containsActionID(response, action.ID)

		if !containsID {
			return false
		}
	}

	return true
}

func containsNoneActions(response []*dto.GetComponentResponse, actions []*orm.Action) bool {
	for _, action := range actions {
		containsID := containsActionID(response, action.ID)

		if containsID {
			return false
		}
	}

	return true
}
