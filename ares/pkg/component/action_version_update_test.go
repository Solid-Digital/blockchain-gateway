package component_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/orm"

	"bitbucket.org/unchain/ares/gen/dto"

	"github.com/Pallinder/go-randomdata"

	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_UpdateActionVersion() {
	params1 := &dto.UpdateComponentVersionRequest{
		Description:   s.helper.StringPtr(randomdata.Paragraph()),
		ExampleConfig: s.helper.StringPtr(randomdata.Paragraph()),
		InputSchema:   []string{"in", "in2", "in3"},
		OutputSchema:  []string{"out", "out2", "out3"},
		Readme:        s.helper.StringPtr(randomdata.Paragraph()),
	}
	actionVersion1, action1, org1, user1 := s.factory.ActionVersionOrgUser(false, true)

	params2 := &dto.UpdateComponentVersionRequest{
		Description:   s.helper.StringPtr(randomdata.Paragraph()),
		ExampleConfig: s.helper.StringPtr(randomdata.Paragraph()),
		InputSchema:   []string{"in", "in2", "in3"},
		OutputSchema:  []string{"out", "out2", "out3"},
		Readme:        s.helper.StringPtr(randomdata.Paragraph()),
		Public:        s.helper.BoolPtr(randomdata.Boolean()),
	}
	actionVersion2, action2, org2, user2 := s.factory.ActionVersionOrgUser(false, true)
	s.helper.MakeSuperAdmin(user2.ID)

	params3 := &dto.UpdateComponentVersionRequest{
		Description:   s.helper.StringPtr(randomdata.Paragraph()),
		ExampleConfig: s.helper.StringPtr(randomdata.Paragraph()),
		InputSchema:   []string{"in", "in2", "in3"},
		OutputSchema:  []string{"out", "out2", "out3"},
		Readme:        s.helper.StringPtr(randomdata.Paragraph()),
		Public:        s.helper.BoolPtr(randomdata.Boolean()),
	}
	actionVersion3, action3, org3, user3 := s.factory.ActionVersionOrgUser(false, true)

	cases := map[string]struct {
		Request       *dto.UpdateComponentVersionRequest
		ActionVersion *orm.ActionVersion
		Action        *orm.Action
		Organization  *orm.Organization
		Principal     *dto.User
		Success       bool
	}{
		"update action version": {
			params1,
			actionVersion1,
			action1,
			org1,
			s.factory.ORMToDTOUser(user1),
			true,
		},

		"set public with admin user": {
			params2,
			actionVersion2,
			action2,
			org2,
			s.factory.ORMToDTOUser(user2),
			true,
		},
		"set public with non-admin user": {
			params3,
			actionVersion3,
			action3,
			org3,
			s.factory.ORMToDTOUser(user3),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.UpdateActionVersion(tc.Request, tc.Organization.Name, tc.Action.Name, tc.ActionVersion.Version, tc.Principal)

			if tc.Success {
				xrequire.NoError(t, err)

				actionVersionFromDB := s.helper.DBGetActionVersion(*response.ID)
				actionFromDB := s.helper.DBGetActionByName(tc.Action.Name)

				require.Equal(t, tc.ActionVersion.Version, *response.Version)
				require.Equal(t, actionFromDB.ID, actionVersionFromDB.ActionID)
				require.Equal(t, *tc.Request.Description, actionFromDB.Description)
				require.Equal(t, *tc.Request.Readme, actionVersionFromDB.Readme)
				require.Equal(t, *tc.Request.ExampleConfig, actionVersionFromDB.ExampleConfig)

				if tc.Request.Public != nil {
					require.Equal(t, *tc.Request.Public, actionVersionFromDB.Public)
					require.Equal(t, *tc.Request.Public, actionFromDB.Public)
					require.Equal(t, *tc.Request.Public, *response.Public)
				}

				s.helper.EqualJSON(tc.Request.InputSchema, actionVersionFromDB.InputSchema)
				s.helper.EqualJSON(tc.Request.OutputSchema, actionVersionFromDB.OutputSchema)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}
