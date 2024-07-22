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

func (s *TestSuite) TestService_CreateActionVersion() {
	action1, org1, user1 := s.factory.ActionOrgUser(false, true)
	params1 := &ares.CreateActionVersionRequest{
		OrgName:      org1.Name,
		Name:         action1.Name,
		Version:      "0.1",
		Principal:    s.factory.ORMToDTOUser(user1),
		ActionFile:   ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
		InputSchema:  []string{"in", "in2", "in3"},
		OutputSchema: []string{"out", "out2", "out3"},
	}

	action2 := s.factory.Action(false, false)
	org2, user2 := s.factory.OrganizationAndUser(true)
	params2 := &ares.CreateActionVersionRequest{
		OrgName:    org2.Name,
		Name:       action2.Name,
		Version:    "0.2",
		Principal:  s.factory.ORMToDTOUser(user2),
		ActionFile: ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
	}

	action3, _, user3 := s.factory.ActionOrgUser(false, true)
	org3 := s.factory.Organization(false)
	params3 := &ares.CreateActionVersionRequest{
		OrgName:    org3.Name,
		Name:       action3.Name,
		Version:    "0.3",
		Principal:  s.factory.ORMToDTOUser(user3),
		ActionFile: ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
	}

	// use of non existing action to force using the principal
	params4 := &ares.CreateActionVersionRequest{
		OrgName:    s.factory.Organization(true).Name,
		Name:       s.factory.Action(false, false).Name,
		Version:    "0.4",
		Principal:  s.factory.DTOUser(true),
		ActionFile: ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
	}

	// use of non existing action to force using the principal
	params5 := &ares.CreateActionVersionRequest{
		OrgName:    s.factory.Organization(true).Name,
		Name:       s.factory.Action(false, false).Name,
		Version:    "0.5",
		Principal:  s.factory.DTOUser(false),
		ActionFile: ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
	}

	action6, org6, user6 := s.factory.ActionOrgUser(false, true)
	s.helper.MakeSuperAdmin(user6.ID)

	params6 := &ares.CreateActionVersionRequest{
		OrgName:     org6.Name,
		Name:        action6.Name,
		Version:     "0.1",
		Principal:   s.factory.ORMToDTOUser(user6),
		ActionFile:  ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
		InputSchema: []string{"in", "in2", "in3"},

		OutputSchema: []string{"out", "out2", "out3"},
		Public:       s.helper.BoolPtr(randomdata.Boolean()),
	}

	action7, org7, user7 := s.factory.ActionOrgUser(false, true)
	params7 := &ares.CreateActionVersionRequest{
		OrgName:     org7.Name,
		Name:        action7.Name,
		Version:     "0.1",
		Principal:   s.factory.ORMToDTOUser(user7),
		ActionFile:  ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
		InputSchema: []string{"in", "in2", "in3"},

		OutputSchema: []string{"out", "out2", "out3"},
		Public:       s.helper.BoolPtr(true),
	}

	cases := map[string]struct {
		Request *ares.CreateActionVersionRequest
		Success bool
	}{
		"create version for existing action": {
			params1,
			true,
		},

		"create version for non existing action": {
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
		},

		// this is quite unusual behaviour
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
			response, err := s.service.CreateActionVersion(tc.Request)

			if tc.Success {
				xrequire.NoError(t, err)

				actionVersionFromDB := s.helper.DBGetActionVersion(*response.ID)
				actionFromDB := s.helper.DBGetActionByName(tc.Request.Name)

				require.Equal(t, tc.Request.Version, *response.Version)
				require.Equal(t, actionFileName(tc.Request), actionVersionFromDB.FileName)
				require.Equal(t, actionFileID(tc.Request), actionVersionFromDB.FileID)
				require.Equal(t, actionFromDB.ID, actionVersionFromDB.ActionID)

				if tc.Request.Public != nil {
					require.Equal(t, *tc.Request.Public, actionVersionFromDB.Public)
					require.Equal(t, *tc.Request.Public, *response.Public)
				}

				s.helper.EqualJSON(tc.Request.InputSchema, actionVersionFromDB.InputSchema)
				s.helper.EqualJSON(tc.Request.OutputSchema, actionVersionFromDB.OutputSchema)
				require.True(t, s.helper.FileExists(actionFileID(tc.Request)))
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
				require.False(t, s.helper.FileExists(actionFileID(tc.Request)))
			}
		})
	}
}

func actionFileName(params *ares.CreateActionVersionRequest) string {
	return component.ActionVersionFileName(params.Name, params.Version, params.OrgName)
}

func actionFileID(params *ares.CreateActionVersionRequest) string {
	fileName := actionFileName(params)

	return component.ActionVersionFileID(params.Name, params.Version, params.OrgName, fileName)
}
