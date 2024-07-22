package organization_test

import (
	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
)

func (s *OrganizationTestSuite) TestService_UpdateOrganization() {
	org := s.factory.Organization(true)
	params := &dto.UpdateOrganizationRequest{
		DisplayName: "some other name",
	}
	response, err := s.service.UpdateOrganization(params, org.Name)

	xrequire.NoError(s.Suite.T(), err)
	s.Suite.Require().Equal(org.ID, response.ID)
	s.Suite.Require().Equal(params.DisplayName, response.DisplayName)

	orgFromDB := s.helper.DBGetOrgByName(org.Name)
	s.Suite.Require().Equal(params.DisplayName, orgFromDB.DisplayName)
}
