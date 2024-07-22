package pipeline_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/pkg/testhelper"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_DeployConfiguration() {
	fileName := testhelper.Randumb("http.endpoint.so")
	s.factory.File("../../test/fixtures/binary/http.endpoint.so", fileName)

	lastRevisionTestOrg, lastRevisionTestUser := s.factory.OrganizationAndUser(true)
	lastRevisionTestEnv := s.factory.Environment(lastRevisionTestOrg, lastRevisionTestUser)
	lastRevisionTestPipeline := s.factory.Pipeline(true, lastRevisionTestOrg, lastRevisionTestUser)
	s.factory.Configuration(true, lastRevisionTestOrg, lastRevisionTestUser, lastRevisionTestPipeline)
	s.factory.Configuration(true, lastRevisionTestOrg, lastRevisionTestUser, lastRevisionTestPipeline)
	s.factory.Configuration(true, lastRevisionTestOrg, lastRevisionTestUser, lastRevisionTestPipeline)
	s.factory.Configuration(true, lastRevisionTestOrg, lastRevisionTestUser, lastRevisionTestPipeline)
	s.factory.Configuration(true, lastRevisionTestOrg, lastRevisionTestUser, lastRevisionTestPipeline)
	s.factory.Configuration(true, lastRevisionTestOrg, lastRevisionTestUser, lastRevisionTestPipeline)
	lastRevisionTestConfig := s.helper.DBLatestConfigurationRevision(lastRevisionTestPipeline)
	s.factory.TriggerConfigurationWithFile(false, true, lastRevisionTestConfig, lastRevisionTestOrg, fileName)
	s.factory.BaseConfiguration(false, true, lastRevisionTestConfig, lastRevisionTestOrg)
	lastRevisionTestParams := &dto.DeployConfigurationRequest{
		ConfigurationRevision: s.helper.Int64Ptr(-1),
		Replicas:              s.helper.Int64Ptr(2),
	}

	org1, user1 := s.factory.OrganizationAndUser(true)
	env1 := s.factory.Environment(org1, user1)
	pipeline1 := s.factory.Pipeline(true, org1, user1)
	config1 := s.factory.Configuration(true, org1, user1, pipeline1)
	s.factory.TriggerConfigurationWithFile(false, true, config1, org1, fileName)
	s.factory.BaseConfiguration(false, true, config1, org1)
	params1 := &dto.DeployConfigurationRequest{
		ConfigurationRevision: s.helper.Int64Ptr(config1.Revision),
		Replicas:              s.helper.Int64Ptr(1),
	}

	org2, user2 := s.factory.OrganizationAndUser(true)
	env2 := s.factory.Environment(org2, user2)
	pipeline2 := s.factory.Pipeline(true, org2, user2)
	config2 := s.factory.Configuration(true, org2, user2, pipeline2)
	s.factory.TriggerConfigurationWithFile(false, true, config2, org2, fileName)
	s.factory.BaseConfiguration(false, true, config2, org2)
	params2 := &dto.DeployConfigurationRequest{
		ConfigurationRevision: s.helper.Int64Ptr(config2.Revision),
		Replicas:              s.helper.Int64Ptr(2),
	}

	org3, user3 := s.factory.OrganizationAndUser(true)
	env3 := s.factory.Environment(org3, user3)
	s.helper.SetMaxReplicas(env3, int64(5))
	pipeline3 := s.factory.Pipeline(true, org3, user3)
	config3 := s.factory.Configuration(true, org3, user3, pipeline3)
	s.factory.TriggerConfigurationWithFile(false, true, config3, org3, fileName)
	s.factory.BaseConfiguration(false, true, config3, org3)
	params3 := &dto.DeployConfigurationRequest{
		ConfigurationRevision: s.helper.Int64Ptr(config3.Revision),
		Replicas:              s.helper.Int64Ptr(6),
	}

	org4, user4 := s.factory.OrganizationAndUser(true)
	env4 := s.factory.Environment(org4, user4)
	s.helper.SetMaxReplicas(env4, int64(10))
	pipeline4 := s.factory.Pipeline(true, org4, user4)
	config4 := s.factory.Configuration(true, org4, user4, pipeline4)
	s.factory.TriggerConfigurationWithFile(false, true, config4, org4, fileName)
	s.factory.BaseConfiguration(false, true, config4, org4)
	params4 := &dto.DeployConfigurationRequest{
		ConfigurationRevision: s.helper.Int64Ptr(config4.Revision),
	}

	cases := map[string]struct {
		Params           *dto.DeployConfigurationRequest
		Environment      *orm.Environment
		Organization     *orm.Organization
		Pipeline         *orm.Pipeline
		Principal        *dto.User
		ExpectedReplicas int64
		ExpectedRevision int64
		Success          bool
	}{
		"deploy config without revision set to -1": {
			lastRevisionTestParams,
			lastRevisionTestEnv,
			lastRevisionTestOrg,
			lastRevisionTestPipeline,
			s.factory.DTOUser(true),
			2,
			lastRevisionTestConfig.Revision,
			true,
		},
		"deploy with number of replicas set": {
			params1,
			env1,
			org1,
			pipeline1,
			s.factory.DTOUser(true),
			*params1.Replicas,
			*params1.ConfigurationRevision,
			true,
		},
		"pipeline is not of organization": {
			params2,
			env2,
			s.factory.Organization(true),
			pipeline2,
			s.factory.DTOUser(true),
			0,
			*params2.ConfigurationRevision,
			false,
		},
		"too many replicas": {
			params3,
			env3,
			org3,
			pipeline3,
			s.factory.DTOUser(true),
			0,
			*params3.ConfigurationRevision,
			false,
		},
		"replicas not set - default to max replicas": {
			params4,
			env4,
			org4,
			pipeline4,
			s.factory.DTOUser(true),
			10,
			*params4.ConfigurationRevision,
			true,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			deployment, err := s.service.DeployConfiguration(tc.Params, tc.Organization.Name, tc.Pipeline.Name, tc.Environment.Name, tc.Principal)
			deploymentFromDB := s.helper.DBGetDeployment(tc.Pipeline.ID, tc.Environment.ID)

			if tc.Success {
				xrequire.NoError(t, err)
				require.NotNil(t, deploymentFromDB)
				require.Equal(t, tc.ExpectedReplicas, deploymentFromDB.Replicas)
				require.Equal(t, tc.ExpectedReplicas, *deployment.DesiredReplicas)
				require.Equal(t, tc.ExpectedRevision, deployment.Configuration.Revision)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, deploymentFromDB)
			}
		})
	}
}

func (s *TestSuite) TestService_DeployConfiguration_Redeploy() {
	fileName := testhelper.Randumb("http.endpoint.so")
	s.factory.File("../../test/fixtures/binary/http.endpoint.so", fileName)

	org, user := s.factory.OrganizationAndUser(true)
	env := s.factory.Environment(org, user)
	pipeline := s.factory.Pipeline(true, org, user)
	config := s.factory.Configuration(true, org, user, pipeline)
	s.factory.TriggerConfigurationWithFile(false, true, config, org, fileName)
	s.factory.BaseConfiguration(false, true, config, org)

	params1 := &dto.DeployConfigurationRequest{
		ConfigurationRevision: s.helper.Int64Ptr(config.Revision),
		Replicas:              s.helper.Int64Ptr(3),
	}

	deployment, err := s.service.DeployConfiguration(params1, org.Name, pipeline.Name, env.Name, s.factory.ORMToDTOUser(user))
	xrequire.NoError(s.Suite.T(), err)

	deploymentFromDB := s.helper.DBGetDeployment(pipeline.ID, env.ID)
	s.Require().NotNil(deploymentFromDB)
	s.Require().Equal(int64(3), deploymentFromDB.Replicas)
	s.Require().Equal(int64(3), *deployment.DesiredReplicas)

	params2 := &dto.DeployConfigurationRequest{
		ConfigurationRevision: s.helper.Int64Ptr(config.Revision),
		Replicas:              s.helper.Int64Ptr(2),
	}

	// redeploying existing pipeline with different config is allowed
	deployment, err = s.service.DeployConfiguration(params2, org.Name, pipeline.Name, env.Name, s.factory.ORMToDTOUser(user))
	xrequire.NoError(s.Suite.T(), err)

	deploymentFromDB = s.helper.DBGetDeployment(pipeline.ID, env.ID)
	s.Require().NotNil(deploymentFromDB)
	s.Require().Equal(int64(2), deploymentFromDB.Replicas)
	s.Require().Equal(int64(2), *deployment.DesiredReplicas)

	params3 := &dto.DeployConfigurationRequest{
		ConfigurationRevision: s.helper.Int64Ptr(config.Revision),
		Replicas:              s.helper.Int64Ptr(0),
	}

	// redeploying existing pipeline with 0 replicas to shut is down is also allowed
	deployment, err = s.service.DeployConfiguration(params3, org.Name, pipeline.Name, env.Name, s.factory.ORMToDTOUser(user))
	xrequire.NoError(s.Suite.T(), err)

	deploymentFromDB = s.helper.DBGetDeployment(pipeline.ID, env.ID)
	s.Require().NotNil(deploymentFromDB)
	s.Require().Equal(int64(0), deploymentFromDB.Replicas)
	s.Require().Equal(int64(0), *deployment.DesiredReplicas)
}
