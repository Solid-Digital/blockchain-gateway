package component_test

import (
	"testing"

	"github.com/Pallinder/go-randomdata"

	"bitbucket.org/unchain/ares/pkg/testhelper"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_CreateBaseVersion() {
	base1, org1, user1 := s.factory.BaseOrgUser(false, true)
	_, org2, user2 := s.factory.BaseOrgUser(false, true)
	base3, _, user3 := s.factory.BaseOrgUser(false, true)
	base4, org4, user4 := s.factory.BaseOrgUser(false, true)
	s.helper.MakeSuperAdmin(user4.ID)

	base5, org5, user5 := s.factory.BaseOrgUser(false, true)

	cases := map[string]struct {
		Request      *dto.CreateBaseVersionRequest
		Base         *orm.Base
		Organization *orm.Organization
		Principal    *dto.User
		Success      bool
	}{
		"create version for existing base": {
			&dto.CreateBaseVersionRequest{Version: testhelper.Randumb("alpha")},
			base1,
			org1,
			s.factory.ORMToDTOUser(user1),
			true,
		},
		"create version for non existing base": {
			&dto.CreateBaseVersionRequest{Version: testhelper.Randumb("alpha")},
			s.factory.Base(false, false),
			org2,
			s.factory.ORMToDTOUser(user2),
			true,
		},
		"organization does not exist": {
			&dto.CreateBaseVersionRequest{Version: testhelper.Randumb("alpha")},
			base3,
			s.factory.Organization(false),
			s.factory.ORMToDTOUser(user3),
			false,
		},
		"user not member of organization": {
			&dto.CreateBaseVersionRequest{Version: testhelper.Randumb("alpha")},
			s.factory.Base(false, false), // use of non existing base to force using the principal
			s.factory.Organization(true),
			s.factory.DTOUser(true),
			true, // this is quite unusual behaviour
		},
		"user does not exist": {
			&dto.CreateBaseVersionRequest{Version: testhelper.Randumb("alpha")},
			s.factory.Base(false, false), // use of non existing base to force using the principal
			s.factory.Organization(true),
			s.factory.DTOUser(false),
			false,
		},
		"set public with admin user": {
			&dto.CreateBaseVersionRequest{Version: testhelper.Randumb("alpha"), Public: s.helper.BoolPtr(randomdata.Boolean())},
			base4,
			org4,
			s.factory.ORMToDTOUser(user4),
			true,
		},
		"set public with non-admin user": {
			&dto.CreateBaseVersionRequest{Version: testhelper.Randumb("alpha"), Public: s.helper.BoolPtr(randomdata.Boolean())},
			base5,
			org5,
			s.factory.ORMToDTOUser(user5),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.CreateBaseVersion(tc.Request, tc.Organization.Name, tc.Base.Name, tc.Principal)

			if tc.Success {
				xrequire.NoError(t, err)
				require.NotNil(t, response)

				baseVersionFromDB := s.helper.DBGetBaseVersion(*response.ID)
				baseFromDB := s.helper.DBGetBaseByName(tc.Base.Name)

				if tc.Request.Public != nil {
					require.Equal(t, *tc.Request.Public, baseVersionFromDB.Public)
					require.Equal(t, *tc.Request.Public, *response.Public)
				}

				require.Equal(t, tc.Request.Version, *response.Version)
				require.Equal(t, tc.Request.Version, baseVersionFromDB.Version)
				require.Equal(t, tc.Base.Name, baseFromDB.Name)
				require.Equal(t, baseFromDB.ID, baseVersionFromDB.BaseID)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}
