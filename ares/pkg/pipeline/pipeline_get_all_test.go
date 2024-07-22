package pipeline_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_GetAllPipelines() {
	org1, user1 := s.factory.OrganizationAndUser(true)
	pipeline1a := s.factory.Pipeline(true, org1, user1)
	pipeline1b := s.factory.Pipeline(true, org1, user1)

	cases := map[string]struct {
		Organization *orm.Organization
		ExpectedIDs  []int64
		Success      bool
	}{
		"organization has no pipelines": {
			s.factory.Organization(true),
			[]int64{},
			true,
		},
		"organization with pipelines": {
			org1,
			[]int64{pipeline1a.ID, pipeline1b.ID},
			true,
		},
		"organization does not exist": {
			s.factory.Organization(false),
			[]int64{},
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetAllPipelines(tc.Organization.Name)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, len(tc.ExpectedIDs), len(response))
				require.True(t, containsAllIDs(response, tc.ExpectedIDs))
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}

func containsID(orgs dto.GetAllPipelinesResponse, ID int64) bool {
	for _, value := range orgs {
		if *value.ID == ID {
			return true
		}
	}

	return false
}

func containsAllIDs(orgs dto.GetAllPipelinesResponse, IDs []int64) bool {
	for _, value := range IDs {
		if containsID(orgs, value) == false {
			return false
		}
	}

	return true
}
