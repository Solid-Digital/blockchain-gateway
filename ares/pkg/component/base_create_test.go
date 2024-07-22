package component_test

import (
	"testing"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper"
	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"
	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_CreateBase() {
	cases := map[string]struct {
		Request      *dto.CreateComponentRequest
		Organization *orm.Organization
		User         *dto.User
		Success      bool
	}{
		"create base": {
			&dto.CreateComponentRequest{Name: testhelper.Randumb(randomdata.SillyName())},
			s.factory.Organization(true),
			s.factory.DTOUser(true),
			true},
		"organization does not exist": {
			&dto.CreateComponentRequest{Name: testhelper.Randumb(randomdata.SillyName())},
			s.factory.Organization(false),
			s.factory.DTOUser(true),
			false},
		"user does not exist": {
			&dto.CreateComponentRequest{Name: testhelper.Randumb(randomdata.SillyName())},
			s.factory.Organization(true),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.CreateBase(tc.Request, tc.Organization.Name, tc.User)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.Request.Name, *response.Name)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}
