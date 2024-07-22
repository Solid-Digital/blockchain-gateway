package component_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
	"github.com/Pallinder/go-randomdata"
)

func (s *TestSuite) TestService_UpdateTrigger() {
	trigger1, org1, user1 := s.factory.TriggerOrgUser(false, true)
	trigger2, _, user2 := s.factory.TriggerOrgUser(false, true)
	_, org3, user3 := s.factory.TriggerOrgUser(false, true)
	trigger4, org4, _ := s.factory.TriggerOrgUser(false, true)
	trigger5, org5, _ := s.factory.TriggerOrgUser(false, true)

	cases := map[string]struct {
		Request      *dto.UpdateComponentRequest
		Trigger      *orm.Trigger
		Organization *orm.Organization
		Principal    *dto.User
		Success      bool
	}{
		"update trigger": {
			&dto.UpdateComponentRequest{DisplayName: testhelper.Randumb(randomdata.SillyName())},
			trigger1,
			org1,
			s.factory.ORMToDTOUser(user1),
			true,
		},
		"organization does not exist": {
			&dto.UpdateComponentRequest{DisplayName: testhelper.Randumb(randomdata.SillyName())},
			trigger2,
			s.factory.Organization(false),
			s.factory.ORMToDTOUser(user2),
			false,
		},
		"trigger does not exist": {
			&dto.UpdateComponentRequest{DisplayName: testhelper.Randumb(randomdata.SillyName())},
			s.factory.Trigger(false, false),
			org3,
			s.factory.ORMToDTOUser(user3),
			false,
		},
		"user does not exist": {
			&dto.UpdateComponentRequest{DisplayName: testhelper.Randumb(randomdata.SillyName())},
			trigger4,
			org4,
			s.factory.DTOUser(false),
			false,
		},
		"user not member of organization": {
			&dto.UpdateComponentRequest{DisplayName: testhelper.Randumb(randomdata.SillyName())},
			trigger5,
			org5,
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.UpdateTrigger(tc.Request, tc.Organization.Name, tc.Trigger.Name, tc.Principal)
			triggerFromDB := s.helper.DBGetTrigger(tc.Trigger.ID)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.Request.DisplayName, *response.DisplayName)
				require.Equal(t, tc.Request.DisplayName, triggerFromDB.DisplayName)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
				if triggerFromDB != nil {
					require.Equal(t, tc.Trigger.DisplayName, triggerFromDB.DisplayName)
				}
			}
		})
	}
}
