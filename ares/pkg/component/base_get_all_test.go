package component_test

import (
	"testing"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_GetAllBases() {
	preExistingPublicBasesCount := len(s.helper.DBPublicBases())

	org1, _ := s.factory.OrganizationAndUser(true)
	var included1 = make([]*orm.Base, 5)
	for i := 0; i < 5; i++ {
		included1[i] = s.factory.BaseForOrg(org1, false, true)
	}

	var excluded1 = make([]*orm.Base, 5)
	for i := 0; i < 5; i++ {
		excluded1[i] = s.factory.Base(false, true)
	}

	org2, _ := s.factory.OrganizationAndUser(true)
	var included2 = make([]*orm.Base, 5)
	for i := 0; i < 5; i++ {
		included2[i] = s.factory.BaseForOrg(org2, false, true)
	}

	cases := map[string]struct {
		Organization     *orm.Organization
		Available        *bool
		Included         []*orm.Base
		Excluded         []*orm.Base
		PublicBaseCount  int
		PrivateBaseCount int
		Success          bool
	}{
		"should return bases from org1": {
			org1,
			nil,
			included1,
			excluded1,
			preExistingPublicBasesCount,
			5,
			true,
		},
		"with available true": {
			org1,
			s.helper.BoolPtr(true),
			included1,
			excluded1,
			preExistingPublicBasesCount,
			5,
			true,
		},
		"with available false": {
			org1,
			s.helper.BoolPtr(false),
			[]*orm.Base{},
			[]*orm.Base{},
			0, // is this the expected behaviour?
			0, // is this the expected behaviour?
			true,
		},
		"should return bases from org2": {
			org2,
			nil,
			included2,
			excluded1,
			preExistingPublicBasesCount,
			5,
			true,
		},
		"organization has no bases": {
			s.factory.Organization(true),
			nil,
			[]*orm.Base{},
			[]*orm.Base{},
			preExistingPublicBasesCount,
			0,
			true,
		},
		"organization does not exist": {
			s.factory.Organization(false),
			nil,
			[]*orm.Base{},
			[]*orm.Base{},
			preExistingPublicBasesCount,
			0,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetAllBases(tc.Organization.Name, tc.Available)

			if tc.Success {
				xrequire.NoError(t, err)
				publicCount, privateCount := countBases(response)
				require.Equal(t, tc.PublicBaseCount, publicCount)
				require.Equal(t, tc.PrivateBaseCount, privateCount)
				require.True(t, containsAllBases(response, tc.Included))
				require.True(t, containsNoneBases(response, tc.Excluded))
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}

func countBases(responses []*dto.GetComponentResponse) (publicCount int, privateCount int) {
	for _, res := range responses {
		if *res.Public {
			publicCount++
		} else {
			privateCount++
		}
	}

	return publicCount, privateCount
}

func containsBaseID(lst []*dto.GetComponentResponse, ID int64) bool {
	for _, e := range lst {
		if *e.ID == ID {
			return true
		}
	}

	return false
}

func containsAllBases(response []*dto.GetComponentResponse, bases []*orm.Base) bool {
	for _, base := range bases {
		containsID := containsBaseID(response, base.ID)

		if !containsID {
			return false
		}
	}

	return true
}

func containsNoneBases(response []*dto.GetComponentResponse, bases []*orm.Base) bool {
	for _, base := range bases {
		containsID := containsBaseID(response, base.ID)

		if containsID {
			return false
		}
	}

	return true
}
