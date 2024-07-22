package pipeline_test

import (
	"fmt"
	"testing"

	"bitbucket.org/unchain/ares/pkg/factory"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/dto"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_DeleteEnvironmentVariables() {
	org1, user1 := s.factory.OrganizationAndUser(true)
	pipeline1 := s.factory.Pipeline(true, org1, user1)
	env1 := s.factory.Environment(org1, user1)
	envVars1 := s.factory.EnvVars(true, user1, org1, pipeline1, env1, factory.BothSecretsAndVars)

	org2, user2 := s.factory.OrganizationAndUser(true)
	pipeline2 := s.factory.Pipeline(true, org2, user2)
	env2 := s.factory.Environment(org2, user2)
	_ = s.factory.EnvVars(true, user2, org2, pipeline2, env2, factory.BothSecretsAndVars)

	env2Other := s.factory.Environment(org2, user2)
	envVars2Other := s.factory.EnvVars(true, user2, org2, pipeline2, env2Other, factory.BothSecretsAndVars)

	org3, user3 := s.factory.OrganizationAndUser(true)
	pipeline3 := s.factory.Pipeline(true, org3, user3)
	env3 := s.factory.Environment(org3, user3)
	envVars3 := s.factory.EnvVars(true, user3, org3, pipeline3, env3, factory.BothSecretsAndVars)

	// ensure org3Other has the same user/env
	org3Other := s.factory.Organization(true)
	s.factory.AddUserToOrg(user3, org3Other)
	s.factory.EnvironmentWithName(env3.Name, org3Other, user3)

	org4, user4 := s.factory.OrganizationAndUser(true)
	pipeline4 := s.factory.Pipeline(true, org4, user4)
	env4 := s.factory.Environment(org4, user4)
	envVars4 := s.factory.EnvVars(true, user4, org4, pipeline4, env4, factory.BothSecretsAndVars)

	// ensure org3Other has the same user/pipeline
	org4Other := s.factory.Organization(true)
	s.factory.AddUserToOrg(user4, org4Other)
	s.factory.PipelineWithName(pipeline4.Name, true, org4Other, user4)

	cases := map[string]struct {
		Organization *orm.Organization
		Pipeline     *orm.Pipeline
		Environment  *orm.Environment
		VarID        int64
		User         *dto.User
		Success      bool
	}{
		"delete environment variable": {
			org1,
			pipeline1,
			env1,
			envVars1[0].ID,
			s.factory.ORMToDTOUser(user1),
			true,
		},
		"fail when the specified variable id is in another environment": {
			org2,
			pipeline2,
			env2,
			envVars2Other[0].ID,
			s.factory.ORMToDTOUser(user2),
			false,
		},
		"fail when the specified pipeline doesn't exist in the specified organization": {
			org3Other,
			pipeline3,
			env3,
			envVars3[0].ID,
			s.factory.ORMToDTOUser(user3),
			false,
		},
		"fail when the specified environment doesn't exist in the specified organization": {
			org4Other,
			pipeline4,
			env4,
			envVars4[0].ID,
			s.factory.ORMToDTOUser(user4),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			err := s.service.DeleteEnvironmentVariable(tc.Organization.Name, tc.Pipeline.Name, tc.Environment.Name, tc.VarID, tc.User)

			if tc.Success {
				xrequire.NoError(t, err)
				require.False(t, s.helper.DBEnvVarExists(tc.VarID))
			} else {
				fmt.Printf("%+v\n", err)
				xrequire.Error(t, err)
				require.True(t, s.helper.DBEnvVarExists(tc.VarID))
			}
		})
	}
}
