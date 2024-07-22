package component_test

import (
	"testing"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/gen/dto"

	"github.com/Pallinder/go-randomdata"

	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_UpdateTriggerVersion() {
	params1 := &dto.UpdateComponentVersionRequest{
		Description:   s.helper.StringPtr(randomdata.Paragraph()),
		ExampleConfig: s.helper.StringPtr(randomdata.Paragraph()),
		InputSchema:   []string{"in", "in2", "in3"},
		OutputSchema:  []string{"out", "out2", "out3"},
		Readme:        s.helper.StringPtr(randomdata.Paragraph()),
	}
	triggerVersion1, trigger1, org1, user1 := s.factory.TriggerVersionOrgUser(false, true)

	params2 := &dto.UpdateComponentVersionRequest{
		Description:   s.helper.StringPtr(randomdata.Paragraph()),
		ExampleConfig: s.helper.StringPtr(randomdata.Paragraph()),
		InputSchema:   []string{"in", "in2", "in3"},
		OutputSchema:  []string{"out", "out2", "out3"},
		Readme:        s.helper.StringPtr(randomdata.Paragraph()),
		Public:        s.helper.BoolPtr(randomdata.Boolean()),
	}
	triggerVersion2, trigger2, org2, user2 := s.factory.TriggerVersionOrgUser(false, true)
	s.helper.MakeSuperAdmin(user2.ID)

	params3 := &dto.UpdateComponentVersionRequest{
		Description:   s.helper.StringPtr(randomdata.Paragraph()),
		ExampleConfig: s.helper.StringPtr(randomdata.Paragraph()),
		InputSchema:   []string{"in", "in2", "in3"},
		OutputSchema:  []string{"out", "out2", "out3"},
		Readme:        s.helper.StringPtr(randomdata.Paragraph()),
		Public:        s.helper.BoolPtr(randomdata.Boolean()),
	}
	triggerVersion3, trigger3, org3, user3 := s.factory.TriggerVersionOrgUser(false, true)

	cases := map[string]struct {
		Request        *dto.UpdateComponentVersionRequest
		TriggerVersion *orm.TriggerVersion
		Trigger        *orm.Trigger
		Organization   *orm.Organization
		Principal      *dto.User
		Success        bool
	}{
		"update trigger version": {
			params1,
			triggerVersion1,
			trigger1,
			org1,
			s.factory.ORMToDTOUser(user1),
			true,
		},

		"set public with admin user": {
			params2,
			triggerVersion2,
			trigger2,
			org2,
			s.factory.ORMToDTOUser(user2),
			true,
		},
		"set public with non-admin user": {
			params3,
			triggerVersion3,
			trigger3,
			org3,
			s.factory.ORMToDTOUser(user3),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.UpdateTriggerVersion(tc.Request, tc.Organization.Name, tc.Trigger.Name, tc.TriggerVersion.Version, tc.Principal)

			if tc.Success {
				xrequire.NoError(t, err)

				triggerVersionFromDB := s.helper.DBGetTriggerVersion(*response.ID)
				triggerFromDB := s.helper.DBGetTriggerByName(tc.Trigger.Name)

				require.Equal(t, tc.TriggerVersion.Version, *response.Version)
				require.Equal(t, triggerFromDB.ID, triggerVersionFromDB.TriggerID)
				require.Equal(t, *tc.Request.Description, triggerFromDB.Description)
				require.Equal(t, *tc.Request.Readme, triggerVersionFromDB.Readme)
				require.Equal(t, *tc.Request.ExampleConfig, triggerVersionFromDB.ExampleConfig)

				if tc.Request.Public != nil {
					require.Equal(t, *tc.Request.Public, triggerVersionFromDB.Public)
					require.Equal(t, *tc.Request.Public, triggerFromDB.Public)
					require.Equal(t, *tc.Request.Public, *response.Public)
				}

				s.helper.EqualJSON(tc.Request.InputSchema, triggerVersionFromDB.InputSchema)
				s.helper.EqualJSON(tc.Request.OutputSchema, triggerVersionFromDB.OutputSchema)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}
