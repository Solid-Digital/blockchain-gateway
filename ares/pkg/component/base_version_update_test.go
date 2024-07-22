package component_test

import (
	"testing"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/dto"

	"github.com/Pallinder/go-randomdata"

	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_UpdateBaseVersion() {
	params1 := &dto.UpdateBaseVersionRequest{
		Description: s.helper.StringPtr(randomdata.Paragraph()),
		Readme:      s.helper.StringPtr(randomdata.Paragraph()),
	}
	baseVersion1, base1, org1, user1 := s.factory.BaseVersionOrgUser(false, true)

	params2 := &dto.UpdateBaseVersionRequest{
		Description: s.helper.StringPtr(randomdata.Paragraph()),
		Readme:      s.helper.StringPtr(randomdata.Paragraph()),
		Public:      s.helper.BoolPtr(randomdata.Boolean()),
	}
	baseVersion2, base2, org2, user2 := s.factory.BaseVersionOrgUser(false, true)
	s.helper.MakeSuperAdmin(user2.ID)

	params3 := &dto.UpdateBaseVersionRequest{
		Description: s.helper.StringPtr(randomdata.Paragraph()),
		Readme:      s.helper.StringPtr(randomdata.Paragraph()),
		Public:      s.helper.BoolPtr(randomdata.Boolean()),
	}
	baseVersion3, base3, org3, user3 := s.factory.BaseVersionOrgUser(false, true)

	cases := map[string]struct {
		Request      *dto.UpdateBaseVersionRequest
		BaseVersion  *orm.BaseVersion
		Base         *orm.Base
		Organization *orm.Organization
		Principal    *dto.User
		Success      bool
	}{
		"update base version": {
			params1,
			baseVersion1,
			base1,
			org1,
			s.factory.ORMToDTOUser(user1),
			true,
		},

		"set public with admin user": {
			params2,
			baseVersion2,
			base2,
			org2,
			s.factory.ORMToDTOUser(user2),
			true,
		},
		"set public with non-admin user": {
			params3,
			baseVersion3,
			base3,
			org3,
			s.factory.ORMToDTOUser(user3),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.UpdateBaseVersion(tc.Request, tc.Organization.Name, tc.Base.Name, tc.BaseVersion.Version, tc.Principal)

			if tc.Success {
				xrequire.NoError(t, err)

				baseVersionFromDB := s.helper.DBGetBaseVersion(*response.ID)
				baseFromDB := s.helper.DBGetBaseByName(tc.Base.Name)

				require.Equal(t, tc.BaseVersion.Version, *response.Version)
				require.Equal(t, baseFromDB.ID, baseVersionFromDB.BaseID)
				require.Equal(t, *tc.Request.Description, baseFromDB.Description)
				require.Equal(t, *tc.Request.Readme, baseVersionFromDB.Readme)

				if tc.Request.Public != nil {
					require.Equal(t, *tc.Request.Public, baseVersionFromDB.Public)
					require.Equal(t, *tc.Request.Public, baseFromDB.Public)
					require.Equal(t, *tc.Request.Public, *response.Public)
				}
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}
