package pipeline_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/dto"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_GetPipeline() {
	org1, user1 := s.factory.OrganizationAndUser(true)
	pipeline1 := s.factory.Pipeline(true, org1, user1)

	org2, user2 := s.factory.OrganizationAndUser(true)
	pipeline2 := s.factory.Pipeline(true, org2, user2)

	org3, user3 := s.factory.OrganizationAndUser(true)
	pipeline3 := s.factory.Pipeline(true, org3, user3)

	org4, user4 := s.factory.OrganizationAndUser(true)
	pipeline4 := s.factory.Pipeline(false, org4, user4)

	org5, user5 := s.factory.OrganizationAndUser(true)
	env5a := s.factory.Environment(org5, user5)
	env5b := s.factory.Environment(org5, user5)
	pipeline5 := s.factory.Pipeline(true, org5, user5)

	org6, user6 := s.factory.OrganizationAndUser(true)
	pipeline6 := s.factory.Pipeline(true, org6, user6)
	configuration6 := s.factory.Configuration(true, org6, user6, pipeline6)
	env6 := s.factory.Environment(org6, user6)
	s.factory.BaseConfiguration(false, true, configuration6, org6)
	s.factory.TriggerConfiguration(false, true, configuration6, org6)

	cases := map[string]struct {
		Organization    *orm.Organization
		User            *orm.User
		Pipeline        *orm.Pipeline
		ExpectedEnvs    []string
		ExpectedConfigs []int64
		Success         bool
	}{
		"get pipeline": {
			org1,
			user1,
			pipeline1,
			[]string{},
			[]int64{},
			true,
		},
		"organization does not exist": {
			s.factory.Organization(false),
			user2,
			pipeline2,
			[]string{},
			[]int64{},
			false,
		},
		"pipeline does not belong to organization": {
			s.factory.Organization(true),
			user3,
			pipeline3,
			[]string{},
			[]int64{},
			false,
		},
		"pipeline does not exist": {
			org4,
			user4,
			pipeline4,
			[]string{},
			[]int64{},
			false,
		},
		"with environments": {
			org5,
			user5,
			pipeline5,
			[]string{env5a.Name, env5b.Name},
			[]int64{},
			true,
		},
		"with configurations": {
			org6,
			user6,
			pipeline6,
			[]string{env6.Name},
			[]int64{configuration6.Revision},
			true,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.GetPipeline(tc.Organization.Name, tc.Pipeline.Name)

			if tc.Success {
				xrequire.NoError(t, err)
				require.NotNil(t, response)
				require.Equal(t, tc.Pipeline.ID, *response.ID)
				require.Equal(t, tc.Pipeline.DisplayName, *response.DisplayName)
				require.Equal(t, tc.Pipeline.Name, *response.Name)
				require.Equal(t, tc.Pipeline.Description, *response.Description)
				require.Equal(t, tc.User.ID, response.CreatedBy.ID)
				require.Equal(t, tc.User.ID, response.UpdatedBy.ID)

				require.Equal(t, len(tc.ExpectedEnvs), len(response.Environments))
				for _, envName := range tc.ExpectedEnvs {
					require.True(t, envSliceContains(response.Environments, envName))
				}

				require.Equal(t, len(tc.ExpectedConfigs), len(response.ConfigurationRevisions))
				for i, cfg := range tc.ExpectedConfigs {
					require.Equal(t, cfg, response.ConfigurationRevisions[i].Revision)
				}
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}

func envSliceContains(envs []*dto.PipelineEnvironment, envName string) bool {
	for _, env := range envs {
		if env.Name == envName {
			return true
		}
	}

	return false
}
