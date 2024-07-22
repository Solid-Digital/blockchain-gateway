package pipeline_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_GetDeployment() {
	org1, pipeline1, env1, deployment1 := s.factory.DeploymentFromService()

	_, pipeline2, env2, deployment2 := s.factory.DeploymentFromService()

	cases := map[string]struct {
		Organization *orm.Organization
		Pipeline     *orm.Pipeline
		Environment  *orm.Environment
		Deployment   *orm.Deployment
		Success      bool
	}{
		"get deployment": {
			org1,
			pipeline1,
			env1,
			deployment1,
			true,
		},
		"pipeline not from organization": {
			s.factory.Organization(true),
			pipeline2,
			env2,
			deployment2,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			deployment, err := s.service.GetDeployment(tc.Organization.Name, tc.Pipeline.Name, tc.Environment.Name)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, *deployment.ID, tc.Deployment.ID)
				require.Equal(t, *deployment.DesiredReplicas, tc.Deployment.Replicas)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, deployment)
				require.True(t, s.helper.DBDeploymentByIdExists(tc.Deployment.ID))
			}
		})
	}
}
