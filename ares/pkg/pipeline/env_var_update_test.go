package pipeline_test

import (
	"fmt"
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/pkg/factory"

	"bitbucket.org/unchain/ares/pkg/testhelper"
	"github.com/Pallinder/go-randomdata"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_UpdateEnvironmentVariables() {
	org1, user1 := s.factory.OrganizationAndUser(true)
	pipeline1 := s.factory.Pipeline(true, org1, user1)
	env1 := s.factory.Environment(org1, user1)
	envVars1 := s.factory.EnvVars(true, user1, org1, pipeline1, env1, factory.BothSecretsAndVars)
	params1 := &dto.UpdateEnvironmentVariablesRequest{
		Key:   testhelper.StringPtr(testhelper.Randumb(randomdata.Noun())),
		Value: testhelper.StringPtr(testhelper.Randumb(randomdata.Noun())),
	}

	org2, user2 := s.factory.OrganizationAndUser(true)
	pipeline2 := s.factory.Pipeline(true, org2, user2)
	env2 := s.factory.Environment(org2, user2)
	envVars2 := s.factory.EnvVars(true, user2, org2, pipeline2, env2, factory.OnlySecrets)
	params2 := &dto.UpdateEnvironmentVariablesRequest{
		Key:   testhelper.StringPtr(testhelper.Randumb(randomdata.Noun())),
		Value: testhelper.StringPtr(testhelper.Randumb(randomdata.Noun())),
	}

	org3, user3 := s.factory.OrganizationAndUser(true)
	pipeline3 := s.factory.Pipeline(true, org3, user3)
	env3 := s.factory.Environment(org3, user3)
	envVars3 := s.factory.EnvVars(true, user3, org3, pipeline3, env3, factory.OnlyVars)
	params3 := &dto.UpdateEnvironmentVariablesRequest{
		Key:   testhelper.StringPtr(testhelper.Randumb(randomdata.Noun())),
		Value: testhelper.StringPtr(testhelper.Randumb(randomdata.Noun())),
	}

	org4, user4 := s.factory.OrganizationAndUser(true)
	pipeline4 := s.factory.Pipeline(true, org4, user4)
	env4 := s.factory.Environment(org4, user4)
	envVars4 := s.factory.EnvVars(true, user4, org4, pipeline4, env4, factory.BothSecretsAndVars)

	params4 := &dto.UpdateEnvironmentVariablesRequest{
		Key:   testhelper.StringPtr(testhelper.Randumb(randomdata.Noun())),
		Value: testhelper.StringPtr(testhelper.Randumb(randomdata.Noun())),
	}

	// ensure org4Other has the same user/env
	org4Other := s.factory.Organization(true)
	s.factory.AddUserToOrg(user4, org4Other)
	s.factory.EnvironmentWithName(env4.Name, org4Other, user4)

	org5, user5 := s.factory.OrganizationAndUser(true)
	pipeline5 := s.factory.Pipeline(true, org5, user5)
	env5 := s.factory.Environment(org5, user5)
	envVars5 := s.factory.EnvVars(true, user5, org5, pipeline5, env5, factory.BothSecretsAndVars)

	params5 := &dto.UpdateEnvironmentVariablesRequest{
		Key:   testhelper.StringPtr(testhelper.Randumb(randomdata.Noun())),
		Value: testhelper.StringPtr(testhelper.Randumb(randomdata.Noun())),
	}

	// ensure org5Other has the same user/pipeline
	org5Other := s.factory.Organization(true)
	s.factory.AddUserToOrg(user5, org5Other)
	s.factory.PipelineWithName(pipeline5.Name, true, org5Other, user5)

	org6, user6 := s.factory.OrganizationAndUser(true)
	pipeline6 := s.factory.Pipeline(true, org6, user6)
	env6 := s.factory.Environment(org6, user6)
	_ = s.factory.EnvVars(true, user6, org6, pipeline6, env6, factory.BothSecretsAndVars)

	env6Other := s.factory.Environment(org6, user6)
	envVars6Other := s.factory.EnvVars(true, user6, org6, pipeline6, env6Other, factory.BothSecretsAndVars)

	params6 := &dto.UpdateEnvironmentVariablesRequest{
		Key:   testhelper.StringPtr(testhelper.Randumb(randomdata.Noun())),
		Value: testhelper.StringPtr(testhelper.Randumb(randomdata.Noun())),
	}

	cases := map[string]struct {
		Params       *dto.UpdateEnvironmentVariablesRequest
		Organization *orm.Organization
		Pipeline     *orm.Pipeline
		Environment  *orm.Environment
		VarID        int64
		User         *dto.User
		Success      bool
	}{
		"update environment variable": {
			params1,
			org1,
			pipeline1,
			env1,
			envVars1[0].ID,
			s.factory.ORMToDTOUser(user1),
			true,
		},
		"update environment variable only variables": {
			params2,
			org2,
			pipeline2,
			env2,
			envVars2[0].ID,
			s.factory.ORMToDTOUser(user2),
			true,
		},
		"update environment variable only secrets": {
			params3,
			org3,
			pipeline3,
			env3,
			envVars3[0].ID,
			s.factory.ORMToDTOUser(user3),
			true,
		},
		"fail when the specified pipeline doesn't exist in the specified organization": {
			params4,
			org4Other,
			pipeline4,
			env4,
			envVars4[0].ID,
			s.factory.ORMToDTOUser(user4),
			false,
		},
		"fail when the specified environment doesn't exist in the specified organization": {
			params5,
			org5Other,
			pipeline5,
			env5,
			envVars5[0].ID,
			s.factory.ORMToDTOUser(user5),
			false,
		},
		"fail when the specified variable id is in another environment": {
			params6,
			org6,
			pipeline6,
			env6,
			envVars6Other[0].ID,
			s.factory.ORMToDTOUser(user6),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			res, err := s.service.UpdateEnvironmentVariable(tc.Params, tc.Organization.Name, tc.Pipeline.Name, tc.Environment.Name, tc.VarID, tc.User)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.VarID, res.ID)
				require.Equal(t, *tc.Params.Key, res.Key)
				if res.Secret {
					require.Empty(t, res.Value)
				} else {
					require.Equal(t, *tc.Params.Value, res.Value)
				}
			} else {
				fmt.Printf("%+v\n", err)
				xrequire.Error(t, err)
			}
		})
	}
}
