package organization_test

import (
	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
)

func (s *OrganizationTestSuite) TestService_GetAllOrganizations() {
	/* This returns all organizations, not just the ones from the user */

	org1 := s.factory.Organization(true)
	org2 := s.factory.Organization(true)
	user := s.factory.DTOUser(false)

	orgs, err := s.service.GetAllOrganizations(user)

	xrequire.NoError(s.Suite.T(), err)
	s.Require().True(containsID(orgs, org1.ID))
	s.Require().True(containsID(orgs, org2.ID))
}

func containsID(lst []*dto.GetOrganizationResponse, ID int64) bool {
	for _, e := range lst {
		if e.ID == ID {
			return true
		}
	}

	return false
}
