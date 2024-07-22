package organization_test

import "bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

func (s *OrganizationTestSuite) TestService_GetOrganization() {
	// This implies that any user can get any organization (i.e. without
	// being a member of an organization
	newOrg := s.factory.Organization(true)
	user := s.factory.DTOUser(false)
	org, err := s.service.GetOrganization(newOrg.Name, user)

	xrequire.NoError(s.Suite.T(), err)
	s.Suite.Require().Equal(newOrg.Name, org.Name)
}

func (s *OrganizationTestSuite) TestService_GetOrganizationInvalidName() {
	user := s.factory.DTOUser(false)
	_, err := s.service.GetOrganization("foobar", user)

	xrequire.Error(s.Suite.T(), err)
}
