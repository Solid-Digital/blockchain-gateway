package organization_test

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *OrganizationTestSuite) TestService_RemoveMember() {
	org, user := s.factory.OrganizationAndUser(true)
	roles := s.factory.SomeRoles()
	err := s.ares.Enforcer.SetMemberRoles(org.Name, user.ID, roles)

	xrequire.NoError(s.Suite.T(), err)

	err = s.service.RemoveMember(user.Email.String, org.Name)

	xrequire.NoError(s.Suite.T(), err)

	var isMember bool
	err = s.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		isMember, err = org.Users(orm.UserWhere.ID.EQ(user.ID)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	xrequire.NoError(s.Suite.T(), err)
	s.Suite.Require().False(isMember)

	userRoles := s.ares.Enforcer.GetRolesForUserInOrganization(user.ID, org.Name)

	s.Suite.Require().Equal(0, len(userRoles))
}

func (s *OrganizationTestSuite) TestService_RemoveMemberNotMember() {
	user := s.factory.User(true)
	org := s.factory.Organization(true)
	err := s.service.RemoveMember(user.Email.String, org.Name)

	xrequire.Error(s.Suite.T(), err)
}
