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

func (s *TestSuite) TestService_UpdateBase() {
	base1, org1, user1 := s.factory.BaseOrgUser(false, true)
	base2, _, user2 := s.factory.BaseOrgUser(false, true)
	_, org3, user3 := s.factory.BaseOrgUser(false, true)
	base4, org4, _ := s.factory.BaseOrgUser(false, true)
	base5, org5, _ := s.factory.BaseOrgUser(false, true)

	cases := map[string]struct {
		Request      *dto.UpdateComponentRequest
		Base         *orm.Base
		Organization *orm.Organization
		Principal    *dto.User
		Success      bool
	}{
		"update base": {
			&dto.UpdateComponentRequest{DisplayName: testhelper.Randumb(randomdata.SillyName())},
			base1,
			org1,
			s.factory.ORMToDTOUser(user1),
			true,
		},
		"organization does not exist": {
			&dto.UpdateComponentRequest{DisplayName: testhelper.Randumb(randomdata.SillyName())},
			base2,
			s.factory.Organization(false),
			s.factory.ORMToDTOUser(user2),
			false,
		},
		"base does not exist": {
			&dto.UpdateComponentRequest{DisplayName: testhelper.Randumb(randomdata.SillyName())},
			s.factory.Base(false, false),
			org3,
			s.factory.ORMToDTOUser(user3),
			false,
		},
		"user does not exist": {
			&dto.UpdateComponentRequest{DisplayName: testhelper.Randumb(randomdata.SillyName())},
			base4,
			org4,
			s.factory.DTOUser(false),
			false,
		},
		"user not member of organization": {
			&dto.UpdateComponentRequest{DisplayName: testhelper.Randumb(randomdata.SillyName())},
			base5,
			org5,
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.UpdateBase(tc.Request, tc.Organization.Name, tc.Base.Name, tc.Principal)
			baseFromDB := s.helper.DBGetBase(tc.Base.ID)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.Request.DisplayName, *response.DisplayName)
				require.Equal(t, tc.Request.DisplayName, baseFromDB.DisplayName)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
				if baseFromDB != nil {
					require.Equal(t, tc.Base.DisplayName, baseFromDB.DisplayName)
				}
			}
		})
	}
}
