package pipeline_test

import (
	"fmt"
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"github.com/stretchr/testify/require"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *TestSuite) TestService_CreatePipeline() {
	org1, user1 := s.factory.OrganizationAndUser(true)
	params1 := &dto.CreatePipelineRequest{
		DisplayName: fmt.Sprintf("%s Pipeline", org1.DisplayName),
		Name:        fmt.Sprintf("%s-pipeline", org1.Name),
		Description: "my pipeline description",
	}

	org2 := s.factory.Organization(true)
	params2 := &dto.CreatePipelineRequest{
		DisplayName: fmt.Sprintf("%s Pipeline", org2.DisplayName),
		Name:        fmt.Sprintf("%s-pipeline", org2.Name),
		Description: "my pipeline description",
	}

	org3 := s.factory.Organization(true)
	params3 := &dto.CreatePipelineRequest{
		DisplayName: fmt.Sprintf("%s Pipeline", org3.DisplayName),
		Name:        fmt.Sprintf("%s-pipeline", org3.Name),
		Description: "my pipeline description",
	}

	org4 := s.factory.Organization(false)
	params4 := &dto.CreatePipelineRequest{
		DisplayName: fmt.Sprintf("%s Pipeline", org4.DisplayName),
		Name:        fmt.Sprintf("%s-pipeline", org4.Name),
		Description: "my pipeline description",
	}

	org5, user5 := s.factory.OrganizationAndUser(true)
	env5 := s.factory.Environment(org5, user5)
	params5 := &dto.CreatePipelineRequest{
		DisplayName: fmt.Sprintf("%s Pipeline", org5.DisplayName),
		Name:        fmt.Sprintf("%s-pipeline", org5.Name),
		Description: "my pipeline description",
	}

	cases := map[string]struct {
		Params       *dto.CreatePipelineRequest
		Organization *orm.Organization
		Principal    *dto.User
		ExpectedEnvs []string
		Success      bool
	}{
		"create pipeline": {
			params1,
			org1,
			s.factory.ORMToDTOUser(user1),
			[]string{},
			true,
		},
		"principal not member of organization": {
			params2,
			org2,
			s.factory.DTOUser(true),
			[]string{},
			true,
		},
		"principal does not exist": {
			params3,
			org3,
			s.factory.DTOUser(false),
			[]string{},
			false,
		},
		"organization does not exist": {
			params4,
			org4,
			s.factory.DTOUser(true),
			[]string{},
			false,
		},
		"with environment": {
			params5,
			org5,
			s.factory.ORMToDTOUser(user5),
			[]string{env5.Name},
			true,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.CreatePipeline(tc.Params, tc.Organization.Name, tc.Principal)

			if tc.Success {
				xrequire.NoError(t, err)
				require.NotNil(t, response)
				require.Equal(t, tc.Params.Name, *response.Name)
				require.Equal(t, tc.Params.DisplayName, *response.DisplayName)
				require.Equal(t, tc.Params.Description, *response.Description)
				require.Equal(t, tc.Principal.ID, response.CreatedBy.ID)
				require.Equal(t, tc.Principal.ID, response.UpdatedBy.ID)

				pipelineFromDB := s.helper.DBGetPipeline(*response.ID)

				require.NotNil(t, pipelineFromDB)
				require.Equal(t, tc.Params.Name, pipelineFromDB.Name)
				require.Equal(t, tc.Params.DisplayName, pipelineFromDB.DisplayName)
				require.Equal(t, tc.Params.Description, pipelineFromDB.Description)
				require.Equal(t, tc.Principal.ID, pipelineFromDB.CreatedByID)
				require.Equal(t, tc.Principal.ID, pipelineFromDB.UpdatedByID)
				require.Equal(t, tc.Organization.ID, pipelineFromDB.OrganizationID)

				draftConfigurationFromDB := s.helper.DBGetDraftConfiguration(s.helper.ToInt64(pipelineFromDB.DraftConfigurationID))

				require.NotNil(t, draftConfigurationFromDB)
				require.Equal(t, tc.Principal.ID, draftConfigurationFromDB.CreatedByID)
				require.Equal(t, tc.Principal.ID, draftConfigurationFromDB.UpdatedByID)
				require.Equal(t, tc.Organization.ID, draftConfigurationFromDB.OrganizationID)
				require.Equal(t, int64(1), draftConfigurationFromDB.Revision)

				require.True(t, s.helper.DBBaseDraftConfigurationByDraftConfigurationIDExists(draftConfigurationFromDB.ID))
				require.True(t, s.helper.DBTriggerDraftConfigurationByDraftConfigurationIDExists(draftConfigurationFromDB.ID))

				require.Equal(t, len(tc.ExpectedEnvs), len(response.Environments))
				for _, envName := range tc.ExpectedEnvs {
					require.True(t, envSliceContains(response.Environments, envName))
				}
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)

				pipelineFromDB := s.helper.DBGetPipelineByName(tc.Params.Name)

				require.Nil(t, pipelineFromDB)
			}
		})
	}
}
