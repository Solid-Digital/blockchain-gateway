package pipeline_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_DeletePipeline() {
	org1, user1 := s.factory.OrganizationAndUser(true)
	pipeline1 := s.factory.Pipeline(true, org1, user1)

	_, pipeline2, _, _ := s.factory.DeploymentFromService()

	org3, pipeline3, _, _ := s.factory.DeploymentFromService()

	cases := map[string]struct {
		Organization *orm.Organization
		Pipeline     *orm.Pipeline
		Principal    *dto.User
		Success      bool
	}{
		"delete pipeline": {
			org1,
			pipeline1,
			s.factory.ORMToDTOUser(user1),
			true,
		},
		"pipeline is of different org": {
			s.factory.Organization(true),
			pipeline2,
			s.factory.ORMToDTOUser(user1),
			false,
		},
		"pipeline has active deployment": {
			org3,
			pipeline3,
			s.factory.DTOUser(true),
			true,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			err := s.service.DeletePipeline(tc.Organization.Name, tc.Pipeline.Name, tc.Principal)

			if tc.Success {
				xrequire.NoError(t, err)
				require.False(t, s.helper.DBPipelineExists(tc.Pipeline.ID))
				require.False(t, s.helper.HasActiveDeployments(tc.Pipeline))
			} else {
				xrequire.Error(t, err)
				require.True(t, s.helper.DBPipelineExists(tc.Pipeline.ID))
				require.True(t, s.helper.HasActiveDeployments(tc.Pipeline))
			}
		})
	}
}
