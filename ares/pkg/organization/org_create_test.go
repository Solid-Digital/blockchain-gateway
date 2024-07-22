package organization_test

import (
	"database/sql"
	"encoding/json"
	"errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"github.com/davecgh/go-spew/spew"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/dto"
)

func (s *OrganizationTestSuite) TestService_CreateOrganization() {
	// Make sure that we have default envs (with non default max replicas)
	s.factory.DefaultEnvironment(true)
	defaultEnvs := s.helper.GetDefaultEnvironments()
	defaultEnvsMap := s.helper.DefaultEnvironmentsToMap(defaultEnvs)

	// This should create a user in the db
	user := s.factory.DTOUser(true)
	newOrg := s.factory.Organization(false)

	params := &dto.CreateOrganizationRequest{
		DisplayName: newOrg.DisplayName,
		Name:        newOrg.Name,
	}

	org, err := s.service.CreateOrganization(params, user)

	xrequire.NoError(s.Suite.T(), err)
	s.Require().Equal(params.DisplayName, org.DisplayName)

	orgFromDB := s.helper.DBGetOrgByID(org.ID)

	xrequire.NoError(s.Suite.T(), err)
	s.Require().Equal(params.DisplayName, orgFromDB.DisplayName)
	s.Require().Equal(len(orgFromDB.R.Environments), len(defaultEnvs))
	for _, env := range orgFromDB.R.Environments {
		defaultEnv := defaultEnvsMap[env.Name]
		s.Require().Equal(defaultEnv.MaxReplicas, env.MaxReplicas)
	}
}

// Test that you cannot create an organization for a non existing user
func (s *OrganizationTestSuite) TestService_CreateOrganizationUserNotExists() {
	// This user does not exist in the db
	user := s.factory.DTOUser(false)
	newOrg := s.factory.Organization(false)

	params := &dto.CreateOrganizationRequest{
		DisplayName: newOrg.DisplayName,
		Name:        newOrg.Name,
	}

	org, err := s.service.CreateOrganization(params, user)
	x, _ := json.Marshal(err)
	spew.Dump(string(x))
	s.Require().True(errors.Is(err, sql.ErrNoRows))
	s.Require().True(errors.Is(err, apperr.NotFound))
	s.Require().True(errors.Is(err, ares.ErrUserIDNotFound(apperr.NotFound.Copy(), user.ID)))
	s.Require().Nil(org)
	s.Require().False(s.helper.DBOrgExists(params.Name))
}

// Test that you cannot create the same organization twice
func (s *OrganizationTestSuite) TestService_CreateOrganizationDuplicate() {
	user := s.factory.DTOUser(false)
	org := s.factory.Organization(true)

	params := &dto.CreateOrganizationRequest{
		DisplayName: org.DisplayName,
		Name:        org.Name,
	}

	newOrg, err := s.service.CreateOrganization(params, user)

	// Do we need to verify that the right thing went wrong here?

	xrequire.Error(s.Suite.T(), err)
	s.Require().Nil(newOrg)
}
