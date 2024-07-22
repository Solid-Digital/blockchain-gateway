package organization_test

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"github.com/go-openapi/strfmt"

	"github.com/volatiletech/null"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

/*
	Since casbin is a separate package, maybe we should not assert any roles at all,
	just check that some role has been set.
*/

func (s *OrganizationTestSuite) TestService_InviteMember() {
	// user already exists, so no invite is sent
	// user is not member of organization
	user := s.factory.User(true)
	org := s.factory.Organization(true)
	roles := s.factory.SomeRoles()
	params := &dto.InviteMemberRequest{
		Email: strfmt.Email(user.Email.String),
		Roles: roles,
	}
	response, appErr := s.service.InviteMember(params, org.Name)

	xrequire.NoError(s.Suite.T(), appErr)
	s.Require().Empty(response.InviteID)
	s.Require().Empty(s.helper.GetMailbox(user.Email.String))

	var isMember bool
	var err error
	err = s.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		isMember, err = org.Users(orm.UserWhere.ID.EQ(user.ID)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	xrequire.NoError(s.Suite.T(), err)
	s.Require().True(isMember)

	memberRoles := s.ares.Enforcer.GetRolesForUserInOrganization(user.ID, org.Name)

	s.Require().Equal(roles, memberRoles)
}

func (s *OrganizationTestSuite) TestService_InviteMemberNewUser() {
	// user does not exist, so invite should be sent, which implies that
	// the user is not added to the organization yet.

	// user is not member of organization
	user := s.factory.User(false)
	org := s.factory.Organization(true)
	roles := s.factory.SomeRoles()
	params := &dto.InviteMemberRequest{
		Email: strfmt.Email(user.Email.String),
		Roles: roles,
	}
	response, appErr := s.service.InviteMember(params, org.Name)

	xrequire.NoError(s.Suite.T(), appErr)
	s.Require().NotEmpty(response.InviteID)

	var acExists bool
	var dbMember *orm.User
	var err error
	err = s.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		dbMember, err = org.Users(orm.UserWhere.Email.EQ(null.StringFrom(response.Email))).One(ctx, tx)
		if err != nil {
			return err
		}

		acExists, err = orm.AccountConfirmationTokens(orm.AccountConfirmationTokenWhere.UserID.EQ(dbMember.ID)).Exists(ctx, tx)
		if err != nil {
			return err
		}
		s.Require().True(acExists)
		return nil
	})

	xrequire.NoError(s.Suite.T(), err)

	memberRoles := s.ares.Enforcer.GetRolesForUserInOrganization(response.ID, org.Name)

	s.Require().Equal(roles, memberRoles)
}

func (s *OrganizationTestSuite) TestService_InviteMemberAlreadyMember() {
	org, user := s.factory.OrganizationAndUser(true)
	params := &dto.InviteMemberRequest{
		Email: strfmt.Email(user.Email.String),
	}
	_, appErr := s.service.InviteMember(params, org.Name)
	xrequire.Error(s.Suite.T(), appErr)
}
