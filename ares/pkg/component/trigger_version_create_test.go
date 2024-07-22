package component_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/Pallinder/go-randomdata"

	"bitbucket.org/unchain/ares/pkg/component"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_CreateTriggerVersion() {
	trigger1, org1, user1 := s.factory.TriggerOrgUser(false, true)
	params1 := &ares.CreateTriggerVersionRequest{
		OrgName:      org1.Name,
		Name:         trigger1.Name,
		Version:      "0.1",
		Principal:    s.factory.ORMToDTOUser(user1),
		TriggerFile:  ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
		InputSchema:  []string{"in", "in2", "in3"},
		OutputSchema: []string{"out", "out2", "out3"},
	}

	trigger2 := s.factory.Trigger(false, false)
	org2, user2 := s.factory.OrganizationAndUser(true)
	params2 := &ares.CreateTriggerVersionRequest{
		OrgName:     org2.Name,
		Name:        trigger2.Name,
		Version:     "0.2",
		Principal:   s.factory.ORMToDTOUser(user2),
		TriggerFile: ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
	}

	trigger3, _, user3 := s.factory.TriggerOrgUser(false, true)
	org3 := s.factory.Organization(false)
	params3 := &ares.CreateTriggerVersionRequest{
		OrgName:     org3.Name,
		Name:        trigger3.Name,
		Version:     "0.3",
		Principal:   s.factory.ORMToDTOUser(user3),
		TriggerFile: ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
	}

	// use of non existing trigger to force using the principal
	params4 := &ares.CreateTriggerVersionRequest{
		OrgName:     s.factory.Organization(true).Name,
		Name:        s.factory.Trigger(false, false).Name,
		Version:     "0.4",
		Principal:   s.factory.DTOUser(true),
		TriggerFile: ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
	}

	// use of non existing trigger to force using the principal
	params5 := &ares.CreateTriggerVersionRequest{
		OrgName:     s.factory.Organization(true).Name,
		Name:        s.factory.Trigger(false, false).Name,
		Version:     "0.5",
		Principal:   s.factory.DTOUser(false),
		TriggerFile: ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
	}

	trigger6, org6, user6 := s.factory.TriggerOrgUser(false, true)
	s.helper.MakeSuperAdmin(user6.ID)

	params6 := &ares.CreateTriggerVersionRequest{
		OrgName:     org6.Name,
		Name:        trigger6.Name,
		Version:     "0.1",
		Principal:   s.factory.ORMToDTOUser(user6),
		TriggerFile: ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
		InputSchema: []string{"in", "in2", "in3"},

		OutputSchema: []string{"out", "out2", "out3"},
		Public:       s.helper.BoolPtr(randomdata.Boolean()),
	}

	trigger7, org7, user7 := s.factory.TriggerOrgUser(false, true)
	params7 := &ares.CreateTriggerVersionRequest{
		OrgName:     org7.Name,
		Name:        trigger7.Name,
		Version:     "0.1",
		Principal:   s.factory.ORMToDTOUser(user7),
		TriggerFile: ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
		InputSchema: []string{"in", "in2", "in3"},

		OutputSchema: []string{"out", "out2", "out3"},
		Public:       s.helper.BoolPtr(true),
	}

	cases := map[string]struct {
		Request *ares.CreateTriggerVersionRequest
		Success bool
	}{
		"create version for existing trigger": {
			params1,
			true,
		},
		"create version for non existing trigger": {
			params2,
			true,
		},
		"organization does not exist": {
			params3,
			false,
		},
		"user not member of organization": {
			params4,
			true,
		}, // this is quite unusual behaviour
		"user does not exist": {
			params5,
			false,
		},
		"set public with admin user": {
			params6,
			true,
		},
		"set public with non-admin user": {
			params7,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.CreateTriggerVersion(tc.Request)

			if tc.Success {
				xrequire.NoError(t, err)
				triggerVersionFromDB := s.helper.DBGetTriggerVersion(*response.ID)
				triggerFromDB := s.helper.DBGetTriggerByName(tc.Request.Name)

				require.Equal(t, tc.Request.Version, *response.Version)
				require.Equal(t, triggerFileName(tc.Request), triggerVersionFromDB.FileName)
				require.Equal(t, triggerFileID(tc.Request), triggerVersionFromDB.FileID)
				require.Equal(t, triggerFromDB.ID, triggerVersionFromDB.TriggerID)

				if tc.Request.Public != nil {
					require.Equal(t, *tc.Request.Public, triggerVersionFromDB.Public)
					require.Equal(t, *tc.Request.Public, *response.Public)
				}

				s.helper.EqualJSON(tc.Request.InputSchema, triggerVersionFromDB.InputSchema)
				s.helper.EqualJSON(tc.Request.OutputSchema, triggerVersionFromDB.OutputSchema)
				require.True(t, s.helper.FileExists(triggerFileID(tc.Request)))
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
				require.False(t, s.helper.FileExists(triggerFileID(tc.Request)))
			}
		})
	}
}

func triggerFileName(params *ares.CreateTriggerVersionRequest) string {
	return component.TriggerVersionFileName(params.Name, params.Version, params.OrgName)
}

func triggerFileID(params *ares.CreateTriggerVersionRequest) string {
	fileName := triggerFileName(params)

	return component.TriggerVersionFileID(params.Name, params.Version, params.OrgName, fileName)
}
