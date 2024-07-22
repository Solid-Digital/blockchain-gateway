package organization_test

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *OrganizationTestSuite) TestService_GetAllMembersEmpty() {
	user := s.factory.DTOUser(false)
	org := s.factory.Organization(true)
	response, err := s.service.GetAllMembers(org.Name, user)

	xrequire.NoError(s.Suite.T(), err)
	s.Suite.Require().Equal(0, len(response))
}

func (s *OrganizationTestSuite) TestService_GetAllMembers() {
	org, user1 := s.factory.OrganizationAndUser(true)
	user2 := s.factory.User(true)
	_ = s.factory.User(true) // This user should not be in the response
	err := s.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		err := org.AddUsers(ctx, tx, false, user2)
		if err != nil {
			return err
		}

		return nil
	})

	xrequire.NoError(s.Suite.T(), err)

	randomUser := s.factory.DTOUser(true)
	users, err := s.service.GetAllMembers(org.Name, randomUser)

	xrequire.NoError(s.Suite.T(), err)

	contains := func(users []*dto.GetMemberResponse, user *orm.User) bool {
		for _, v := range users {
			if v.ID == user.ID {
				return true
			}
		}

		return false
	}

	s.Suite.Require().Equal(2, len(users))
	s.Suite.Require().True(contains(users, user1))
	s.Suite.Require().True(contains(users, user2))
}
