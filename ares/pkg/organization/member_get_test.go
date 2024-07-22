package organization_test

import "bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

func (s *OrganizationTestSuite) TestService_GetMember() {
	org, user := s.factory.OrganizationAndUser(true)
	response, err := s.service.GetMember(user.Email.String, org.Name, s.factory.ORMToDTOUser(user))

	xrequire.NoError(s.Suite.T(), err)
	s.Suite.Require().Equal(user.ID, response.ID)
}

func (s *OrganizationTestSuite) TestService_GetMemberInvalidOrg() {
	user := s.factory.DTOUser(true)
	_, err := s.service.GetMember(string(user.Email), "some-org", user)

	xrequire.Error(s.Suite.T(), err)
}

func (s *OrganizationTestSuite) TestService_GetMemberInvalidEmail() {
	org, user := s.factory.OrganizationAndUser(true)
	_, err := s.service.GetMember("some@email.com", org.Name, s.factory.ORMToDTOUser(user))

	xrequire.Error(s.Suite.T(), err)
}

func (s *OrganizationTestSuite) TestService_GetMemberNotOrgMember() {
	org := s.factory.Organization(true)
	user := s.factory.DTOUser(true)
	_, err := s.service.GetMember(string(user.Email), org.Name, user)

	xrequire.Error(s.Suite.T(), err)
}
