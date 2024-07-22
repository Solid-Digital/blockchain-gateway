package http_test

import (
	"math/rand"
	stdhttp "net/http"
	"testing"

	"bitbucket.org/unchain/ares/gen/api/operations/organization"
	"bitbucket.org/unchain/ares/gen/dto"
	mock_ares "bitbucket.org/unchain/ares/gen/mocks"
	"bitbucket.org/unchain/ares/pkg/factory"
	"bitbucket.org/unchain/ares/pkg/http"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/unchainio/pkg/errors"
)

type OrganizationTestSuite struct {
	suite.Suite
	factory       *factory.Factory
	handler       *http.OrganizationHandler
	mockedService *mock_ares.MockOrganizationService
}

func (s *OrganizationTestSuite) SetupSuite() {
	s.T().Skip()
	s.factory = factory.NewFactory(&s.Suite)
}

// This runs before each test
// A new mock needs to be created for each test, otherwise it will fail
func (s *OrganizationTestSuite) SetupTest() {
	s.T().Skip()
	service := mock_ares.NewMockOrganizationService(gomock.NewController(s.T()))
	s.mockedService = service
	s.handler = http.NewOrganizationHandler(service, s.factory.Logger())
}

func (s *OrganizationTestSuite) TestOrganizationHandler_CreateOrganization() {
	cases := map[string]struct {
		ServiceReturn *dto.GetOrganizationResponse
		ServiceError  error
		Params        organization.CreateOrganizationParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetOrganizationResponse{ID: rand.Int63()},
			nil,
			organization.NewCreateOrganizationParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			organization.NewCreateOrganizationParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().CreateOrganization(gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.CreateOrganization(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*organization.CreateOrganizationCreated)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*organization.CreateOrganizationInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *OrganizationTestSuite) TestOrganizationHandler_GetAllMembers() {
	cases := map[string]struct {
		ServiceReturn []*dto.GetMemberResponse
		ServiceError  error
		Params        organization.GetAllMembersParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			[]*dto.GetMemberResponse{{ID: rand.Int63()}},
			nil,
			organization.NewGetAllMembersParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			organization.NewGetAllMembersParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetAllMembers(gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetAllMembers(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*organization.GetAllMembersOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*organization.GetAllMembersInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *OrganizationTestSuite) TestOrganizationHandler_GetAllOrganizations() {
	cases := map[string]struct {
		ServiceReturn []*dto.GetOrganizationResponse
		ServiceError  error
		Params        organization.GetAllOrganizationsParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			[]*dto.GetOrganizationResponse{{ID: rand.Int63()}},
			nil,
			organization.NewGetAllOrganizationsParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			organization.NewGetAllOrganizationsParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetAllOrganizations(gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetAllOrganizations(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*organization.GetAllOrganizationsOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*organization.GetAllOrganizationsInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *OrganizationTestSuite) TestOrganizationHandler_GetMember() {
	cases := map[string]struct {
		ServiceReturn *dto.GetMemberResponse
		ServiceError  error
		Params        organization.GetMemberParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetMemberResponse{ID: rand.Int63()},
			nil,
			organization.NewGetMemberParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			organization.NewGetMemberParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetMember(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetMember(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*organization.GetMemberOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*organization.GetMemberInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *OrganizationTestSuite) TestOrganizationHandler_GetOrganization() {
	cases := map[string]struct {
		ServiceReturn *dto.GetOrganizationResponse
		ServiceError  error
		Params        organization.GetOrganizationParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetOrganizationResponse{ID: rand.Int63()},
			nil,
			organization.NewGetOrganizationParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			organization.NewGetOrganizationParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetOrganization(gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetOrganization(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*organization.GetOrganizationOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*organization.GetOrganizationInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *OrganizationTestSuite) TestOrganizationHandler_InviteMember() {
	cases := map[string]struct {
		ServiceReturn *dto.InviteMemberResponse
		ServiceError  error
		Params        organization.InviteMemberParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.InviteMemberResponse{ID: rand.Int63()},
			nil,
			organization.InviteMemberParams{
				HTTPRequest: &stdhttp.Request{},
			},
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			organization.InviteMemberParams{
				HTTPRequest: &stdhttp.Request{},
			},
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().InviteMember(gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.InviteMember(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*organization.InviteMemberOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*organization.InviteMemberInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *OrganizationTestSuite) TestOrganizationHandler_RemoveMember() {
	cases := map[string]struct {
		ServiceError error
		Params       organization.RemoveMemberParams
		Principal    *dto.User
		Success      bool
	}{
		"service returns no error": {
			nil,
			organization.NewRemoveMemberParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			errors.New("failed"),
			organization.NewRemoveMemberParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().RemoveMember(gomock.Any(), gomock.Any()).Return(tc.ServiceError)
			response := s.handler.RemoveMember(tc.Params, tc.Principal)

			if tc.Success {
				_, ok := response.(*organization.RemoveMemberNoContent)
				require.True(t, ok)
			} else {
				result, ok := response.(*organization.RemoveMemberInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *OrganizationTestSuite) TestOrganizationHandler_SetMemberRoles() {
	cases := map[string]struct {
		ServiceError error
		Params       organization.SetMemberRolesParams
		Principal    *dto.User
		Success      bool
	}{
		"service returns no error": {
			nil,
			organization.NewSetMemberRolesParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			errors.New("failed"),
			organization.NewSetMemberRolesParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().SetMemberRoles(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceError)
			response := s.handler.SetMemberRoles(tc.Params, tc.Principal)

			if tc.Success {
				_, ok := response.(*organization.SetMemberRolesOK)
				require.True(t, ok)
			} else {
				result, ok := response.(*organization.SetMemberRolesInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *OrganizationTestSuite) TestOrganizationHandler_UpdateOrganization() {
	cases := map[string]struct {
		ServiceReturn *dto.GetOrganizationResponse
		ServiceError  error
		Params        organization.UpdateOrganizationParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetOrganizationResponse{ID: rand.Int63()},
			nil,
			organization.NewUpdateOrganizationParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			organization.NewUpdateOrganizationParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().UpdateOrganization(gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.UpdateOrganization(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*organization.UpdateOrganizationOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*organization.UpdateOrganizationInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

// This will make sure the test suite will run
// Don't put any logic for setting up the tests in here, use the hooks from the test suite for that
func TestOrganizationTestSuite(t *testing.T) {
	suite.Run(t, new(OrganizationTestSuite))
}
