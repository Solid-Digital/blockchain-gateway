package pipeline_test

import (
	"fmt"
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/pkg/factory"

	"bitbucket.org/unchain/ares/gen/dto"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_GetAllEnvironmentVariables() {
	org1, user1 := s.factory.OrganizationAndUser(true)
	pipeline1 := s.factory.Pipeline(true, org1, user1)
	env1 := s.factory.Environment(org1, user1)
	envVars1 := s.factory.EnvVars(true, user1, org1, pipeline1, env1, factory.BothSecretsAndVars)

	org2, user2 := s.factory.OrganizationAndUser(true)
	pipeline2 := s.factory.Pipeline(true, org2, user2)
	env2 := s.factory.Environment(org2, user2)
	envVars2 := s.factory.EnvVars(true, user2, org2, pipeline2, env2, factory.BothSecretsAndVars)

	// ensure org2Other has the same user/env
	org2Other := s.factory.Organization(true)
	s.factory.AddUserToOrg(user2, org2Other)
	s.factory.EnvironmentWithName(env2.Name, org2Other, user2)

	org3, user3 := s.factory.OrganizationAndUser(true)
	pipeline3 := s.factory.Pipeline(true, org3, user3)
	env3 := s.factory.Environment(org3, user3)
	envVars3 := s.factory.EnvVars(true, user3, org3, pipeline3, env3, factory.BothSecretsAndVars)

	// ensure org3Other has the same user/pipeline
	org3Other := s.factory.Organization(true)
	s.factory.AddUserToOrg(user3, org3Other)
	s.factory.PipelineWithName(pipeline3.Name, true, org3Other, user3)

	org4, user4 := s.factory.OrganizationAndUser(true)
	pipeline4 := s.factory.Pipeline(true, org4, user4)
	env4 := s.factory.Environment(org4, user4)
	_ = s.factory.EnvVars(true, user4, org4, pipeline4, env4, factory.BothSecretsAndVars)

	org4Other, user4Other := s.factory.OrganizationAndUser(true)
	pipeline4Other := s.factory.Pipeline(true, org4Other, user4Other)
	env4Other := s.factory.Environment(org4Other, user4Other)

	org5, user5 := s.factory.OrganizationAndUser(true)
	pipeline5 := s.factory.Pipeline(true, org5, user5)
	env5 := s.factory.Environment(org5, user5)
	_ = s.factory.EnvVars(true, user5, org5, pipeline5, env5, factory.BothSecretsAndVars)

	org5Other := s.factory.Organization(true)
	s.factory.AddUserToOrg(user5, org5Other)
	pipeline5Other := s.factory.PipelineWithName(pipeline5.Name, true, org5Other, user5)
	env5Other := s.factory.EnvironmentWithName(env5.Name, org5Other, user5)

	org6, user6 := s.factory.OrganizationAndUser(true)
	pipeline6 := s.factory.Pipeline(true, org6, user6)
	env6 := s.factory.Environment(org6, user6)
	_ = s.factory.EnvVars(true, user6, org6, pipeline6, env6, factory.BothSecretsAndVars)

	cases := map[string]struct {
		Organization *orm.Organization
		Pipeline     *orm.Pipeline
		Environment  *orm.Environment
		User         *dto.User
		EnvVars      orm.EnvironmentVariableSlice
		Success      bool
	}{
		"get environment variables": {
			org1,
			pipeline1,
			env1,
			s.factory.ORMToDTOUser(user1),
			envVars1,
			true,
		},
		"fail when the specified pipeline doesn't exist in the specified organization": {
			org2Other,
			pipeline2,
			env2,
			s.factory.ORMToDTOUser(user2),
			envVars2,
			false,
		},
		"fail when the specified environment doesn't exist in the specified organization": {
			org3Other,
			pipeline3,
			env3,
			s.factory.ORMToDTOUser(user3),
			envVars3,
			false,
		},
		"variables added to one organization + pipeline + environment not leaking to another organization + pipeline + environment": {
			org4Other,
			pipeline4Other,
			env4Other,
			s.factory.ORMToDTOUser(user4),
			nil,
			true,
		},
		"variables added to one organization + pipeline + environment not leaking to another organization that has the same pipeline name and the same environment name": {
			org5Other,
			pipeline5Other,
			env5Other,
			s.factory.ORMToDTOUser(user5),
			nil,
			true,
		},
		"variables added to one environment not leaking to another in the same organization + pipeline": {
			org6,
			pipeline6,
			s.factory.Environment(org6, user6),
			s.factory.ORMToDTOUser(user6),
			nil,
			true,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			vars, err := s.service.GetAllEnvironmentVariables(tc.Organization.Name, tc.Pipeline.Name, tc.Environment.Name, tc.User)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Condition(t, func() (success bool) {
					return len(vars) == len(tc.EnvVars)
				})

				for i, v := range vars {
					require.NotEmpty(t, v.Key)
					require.Equal(t, tc.EnvVars[i].Key, v.Key)

					require.Equal(t, tc.EnvVars[i].Secret, v.Secret)
					if v.Secret {
						require.Empty(t, v.Value)
					} else {

						require.NotEmpty(t, v.Value)
						require.Equal(t, tc.EnvVars[i].Value, v.Value)
					}
				}
			} else {
				fmt.Printf("%+v\n", err)
				xrequire.Error(t, err)
				require.Nil(t, vars)
			}
		})
	}
}
