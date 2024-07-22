package pipeline_test

import (
	"fmt"
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/pkg/testhelper"
	"github.com/Pallinder/go-randomdata"

	"bitbucket.org/unchain/ares/gen/dto"

	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_CreateEnvironmentVariables() {
	params1 := &dto.CreateEnvironmentVariableRequest{
		Key:    testhelper.Randumb(randomdata.Noun()),
		Secret: true,
		Value:  testhelper.Randumb(randomdata.Noun()),
	}

	org1, user1 := s.factory.OrganizationAndUser(true)
	pipeline1 := s.factory.Pipeline(true, org1, user1)
	env1 := s.factory.Environment(org1, user1)

	params2 := &dto.CreateEnvironmentVariableRequest{
		Key:    testhelper.Randumb(randomdata.Noun()),
		Secret: false,
		Value:  testhelper.Randumb(randomdata.Noun()),
	}

	org2, user2 := s.factory.OrganizationAndUser(true)
	pipeline2 := s.factory.Pipeline(true, org2, user2)
	env2 := s.factory.Environment(org2, user2)

	params3 := &dto.CreateEnvironmentVariableRequest{
		Key:    testhelper.Randumb(randomdata.Noun()),
		Secret: randomdata.Boolean(),
		Value:  testhelper.Randumb(randomdata.Noun()),
	}

	org3, user3 := s.factory.OrganizationAndUser(true)
	pipeline3 := s.factory.Pipeline(true, org3, user3)
	env3 := s.factory.Environment(org3, user3)

	// ensure org3Other has the same user/env
	org3Other := s.factory.Organization(true)
	s.factory.AddUserToOrg(user3, org3Other)
	s.factory.EnvironmentWithName(env3.Name, org3Other, user3)

	params4 := &dto.CreateEnvironmentVariableRequest{
		Key:    testhelper.Randumb(randomdata.Noun()),
		Secret: randomdata.Boolean(),
		Value:  testhelper.Randumb(randomdata.Noun()),
	}

	org4, user4 := s.factory.OrganizationAndUser(true)
	pipeline4 := s.factory.Pipeline(true, org4, user4)
	env4 := s.factory.Environment(org4, user4)

	// ensure org3Other has the same user/pipeline
	org4Other := s.factory.Organization(true)
	s.factory.AddUserToOrg(user4, org4Other)
	s.factory.PipelineWithName(pipeline4.Name, true, org4Other, user4)

	cases := map[string]struct {
		Params       *dto.CreateEnvironmentVariableRequest
		Organization *orm.Organization
		Pipeline     *orm.Pipeline
		Environment  *orm.Environment
		User         *dto.User
		Success      bool
	}{
		"create environment variable secret": {
			params1,
			org1,
			pipeline1,
			env1,
			s.factory.ORMToDTOUser(user1),
			true,
		},
		"create environment variable non-secret": {
			params2,
			org2,
			pipeline2,
			env2,
			s.factory.ORMToDTOUser(user2),
			true,
		},
		"fail when the specified pipeline doesn't exist in the specified organization": {
			params3,
			org3Other,
			pipeline3,
			env3,
			s.factory.ORMToDTOUser(user3),
			false,
		},
		"fail when the specified environment doesn't exist in the specified organization": {
			params4,
			org4Other,
			pipeline4,
			env4,
			s.factory.ORMToDTOUser(user4),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			res, err := s.service.CreateEnvironmentVariable(tc.Params, tc.Organization.Name, tc.Pipeline.Name, tc.Environment.Name, tc.User)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.Params.Key, res.Key)
				require.Equal(t, tc.Params.Secret, res.Secret)
				if tc.Params.Secret {
					require.Empty(t, res.Value)
				} else {
					require.Equal(t, tc.Params.Value, res.Value)
				}
				require.Equal(t, false, res.Deployed)
				require.Equal(t, tc.User.ID, res.CreatedBy.ID)
				require.Equal(t, tc.User.ID, res.UpdatedBy.ID)
			} else {
				fmt.Printf("%+v\n", err)
				xrequire.Error(t, err)
				require.Nil(t, res)
			}
		})
	}
}
