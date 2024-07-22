package pipeline_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_UpdatePipeline() {
	org1, user1 := s.factory.OrganizationAndUser(true)
	pipeline1 := s.factory.Pipeline(true, org1, user1)

	org2, user2 := s.factory.OrganizationAndUser(true)
	pipeline2 := s.factory.Pipeline(true, org2, user2)

	org3, user3 := s.factory.OrganizationAndUser(true)
	pipeline3 := s.factory.Pipeline(true, org3, user3)

	cases := map[string]struct {
		Params       *dto.UpdatePipelineRequest
		Organization *orm.Organization
		Pipeline     *orm.Pipeline
		Principal    *dto.User
		Success      bool
	}{
		"update pipeline": {
			&dto.UpdatePipelineRequest{
				DisplayName: testhelper.Randumb("updated pipeline"),
			},
			org1,
			pipeline1,
			s.factory.DTOUser(true),
			true,
		},
		"pipeline not from organization": {
			&dto.UpdatePipelineRequest{
				DisplayName: testhelper.Randumb("updated pipeline"),
			},
			s.factory.Organization(true),
			pipeline2,
			s.factory.DTOUser(true),
			false,
		},
		"principal does not exist": {
			&dto.UpdatePipelineRequest{
				DisplayName: testhelper.Randumb("updated pipeline"),
			},
			org3,
			pipeline3,
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.UpdatePipeline(tc.Params, tc.Organization.Name, tc.Pipeline.Name, tc.Principal)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.Params.DisplayName, *response.DisplayName)
				require.Equal(t, tc.Principal.ID, response.UpdatedBy.ID)

				pipelineFromDB := s.helper.DBGetPipeline(tc.Pipeline.ID)

				require.NotNil(t, pipelineFromDB)
				require.Equal(t, tc.Params.DisplayName, pipelineFromDB.DisplayName)
				require.Equal(t, tc.Principal.ID, pipelineFromDB.UpdatedByID)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)

				pipelineFromDB := s.helper.DBGetPipeline(tc.Pipeline.ID)

				require.NotNil(t, pipelineFromDB)
				require.NotEqual(t, tc.Params.DisplayName, pipelineFromDB.DisplayName)
				require.NotEqual(t, tc.Principal.ID, pipelineFromDB.UpdatedByID)
			}
		})
	}
}
