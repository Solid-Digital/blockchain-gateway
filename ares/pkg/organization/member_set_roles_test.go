package organization_test

import (
	"fmt"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/Pallinder/go-randomdata"
)

/*
	Since casbin is a separate package, maybe we should not assert any roles at all,
	just check that some role has been set.
*/

func (s *OrganizationTestSuite) TestService_SetMemberRoles() {
	callingUser := s.factory.DTOUser(false)
	org, user := s.factory.OrganizationAndUser(true)

	params := &dto.SetMemberRolesRequest{
		Roles: map[string]bool{},
	}
	err := s.service.SetMemberRoles(params, user.Email.String, org.Name, callingUser)

	xrequire.NoError(s.Suite.T(), err)

	// The member role is always set in the service
	roles := s.ares.Enforcer.GetRolesForUserInOrganization(user.ID, org.Name)

	s.Suite.Require().True(roles[ares.RoleMember.String()])
}

func (s *OrganizationTestSuite) TestService_SetMemberRolesAllRoles() {
	callingUser := s.factory.DTOUser(false)
	org, user := s.factory.OrganizationAndUser(true)
	roles := s.factory.AllRoles()

	params := &dto.SetMemberRolesRequest{
		Roles: roles,
	}
	err := s.service.SetMemberRoles(params, user.Email.String, org.Name, callingUser)

	xrequire.NoError(s.Suite.T(), err)

	// The member role is always set in the service
	userRoles := s.ares.Enforcer.GetRolesForUserInOrganization(user.ID, org.Name)

	s.Suite.Require().Equal(roles, userRoles)
}

func (s *OrganizationTestSuite) TestService_Randomdata() {
	s.T().Skip()

	data := make(map[string]int)
	for i := 1; i < 10000; i++ {
		rd := randomdata.Email()
		fmt.Println(rd)

		s.Suite.Require().Equal(0, data[rd], "%d: %s", i, rd)
		data[rd] = i
	}
}

func (s *OrganizationTestSuite) TestService_SetMemberRolesRoleFalse() {
	// Assert that roles set to false are not set
	callingUser := s.factory.DTOUser(false)
	org, user := s.factory.OrganizationAndUser(true)

	params := &dto.SetMemberRolesRequest{
		Roles: map[string]bool{ares.RoleUserAdmin.String(): false},
	}
	err := s.service.SetMemberRoles(params, user.Email.String, org.Name, callingUser)

	xrequire.NoError(s.Suite.T(), err)

	// The member role is always set in the service
	roles := s.ares.Enforcer.GetRolesForUserInOrganization(user.ID, org.Name)

	s.Suite.Require().False(roles[ares.RoleUserAdmin.String()])
}
