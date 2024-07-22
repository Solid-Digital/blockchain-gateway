package component_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"github.com/stretchr/testify/require"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper"
	"github.com/Pallinder/go-randomdata"
)

func (s *TestSuite) TestService_UpdateAction() {
	action1, org1, user1 := s.factory.ActionOrgUser(false, true)
	action2, _, user2 := s.factory.ActionOrgUser(false, true)
	_, org3, user3 := s.factory.ActionOrgUser(false, true)
	action4, org4, _ := s.factory.ActionOrgUser(false, true)
	action5, org5, _ := s.factory.ActionOrgUser(false, true)

	cases := map[string]struct {
		Request      *dto.UpdateComponentRequest
		Action       *orm.Action
		Organization *orm.Organization
		Principal    *dto.User
		Success      bool
	}{
		"update action": {
			&dto.UpdateComponentRequest{DisplayName: testhelper.Randumb(randomdata.SillyName())},
			action1,
			org1,
			s.factory.ORMToDTOUser(user1),
			true,
		},
		"organization does not exist": {
			&dto.UpdateComponentRequest{DisplayName: testhelper.Randumb(randomdata.SillyName())},
			action2,
			s.factory.Organization(false),
			s.factory.ORMToDTOUser(user2),
			false,
		},
		"action does not exist": {
			&dto.UpdateComponentRequest{DisplayName: testhelper.Randumb(randomdata.SillyName())},
			s.factory.Action(false, false),
			org3,
			s.factory.ORMToDTOUser(user3),
			false,
		},
		"user does not exist": {
			&dto.UpdateComponentRequest{DisplayName: testhelper.Randumb(randomdata.SillyName())},
			action4,
			org4,
			s.factory.DTOUser(false),
			false,
		},
		"user not member of organization": {
			&dto.UpdateComponentRequest{DisplayName: testhelper.Randumb(randomdata.SillyName())},
			action5,
			org5,
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.UpdateAction(tc.Request, tc.Organization.Name, tc.Action.Name, tc.Principal)
			actionFromDB := s.helper.DBGetAction(tc.Action.ID)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.Request.DisplayName, *response.DisplayName)
				require.Equal(t, tc.Request.DisplayName, actionFromDB.DisplayName)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
				if actionFromDB != nil {
					require.Equal(t, tc.Action.DisplayName, actionFromDB.DisplayName)
				}
			}
		})
	}
}
