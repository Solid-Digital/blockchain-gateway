package http_test

import (
	"context"
	"math/rand"
	stdhttp "net/http"
	"testing"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"github.com/go-chi/chi/middleware"

	"bitbucket.org/unchain/ares/pkg/testhelper"

	"bitbucket.org/unchain/ares/gen/api/operations/component"
	"bitbucket.org/unchain/ares/gen/dto"
	"github.com/stretchr/testify/require"
	"github.com/unchainio/pkg/errors"

	mock_ares "bitbucket.org/unchain/ares/gen/mocks"
	"bitbucket.org/unchain/ares/pkg/factory"
	"bitbucket.org/unchain/ares/pkg/http"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ComponentTestSuite struct {
	suite.Suite
	factory       *factory.Factory
	handler       *http.ComponentHandler
	mockedService *mock_ares.MockComponentService
}

func (s *ComponentTestSuite) SetupSuite() {
	s.factory = factory.NewFactory(&s.Suite)
}

// This runs before each test
// A new mock needs to be created for each test, otherwise it will fail
func (s *ComponentTestSuite) SetupTest() {
	s.T().Skip()
	service := mock_ares.NewMockComponentService(gomock.NewController(s.T()))
	s.mockedService = service
	s.handler = http.NewComponentHandler(service, s.factory.Logger())
}

func NewReq() *stdhttp.Request {
	req := &stdhttp.Request{Header: stdhttp.Header{}}
	req = req.WithContext(context.WithValue(context.Background(), middleware.RequestIDKey, "requestID"))

	return req
}

func (s *ComponentTestSuite) TestComponentHandler_CreateAction() {
	cases := map[string]struct {
		ServiceReturn *dto.GetComponentResponse
		ServiceError  *apperr.Error
		Params        component.CreateActionParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetComponentResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			component.CreateActionParams{
				HTTPRequest: NewReq(),
			},
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			nil,
			apperr.Internal.Copy(),
			component.CreateActionParams{
				HTTPRequest: NewReq(),
			},
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().CreateAction(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)

			response := s.handler.CreateAction(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.CreateActionOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*apperr.Error)
				require.True(t, ok)
				require.Equal(t, result.Status, int64(stdhttp.StatusInternalServerError))
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_CreateActionVersion() {
	cases := map[string]struct {
		ServiceReturn *dto.GetComponentVersionResponse
		ServiceError  error
		Params        component.CreateActionVersionParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetComponentVersionResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			component.NewCreateActionVersionParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewCreateActionVersionParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().CreateActionVersion(gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.CreateActionVersion(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.CreateActionVersionOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.CreateActionVersionInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_CreateBase() {
	cases := map[string]struct {
		ServiceReturn *dto.GetComponentResponse
		ServiceError  error
		Params        component.CreateBaseParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetComponentResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			component.NewCreateBaseParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewCreateBaseParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().CreateBase(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.CreateBase(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.CreateBaseOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.CreateBaseInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_CreateBaseVersion() {
	cases := map[string]struct {
		ServiceReturn *dto.GetBaseVersionResponse
		ServiceError  error
		Params        component.CreateBaseVersionParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetBaseVersionResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			component.NewCreateBaseVersionParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewCreateBaseVersionParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().CreateBaseVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.CreateBaseVersion(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.CreateBaseVersionOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.CreateBaseVersionInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_CreateTrigger() {
	cases := map[string]struct {
		ServiceReturn *dto.GetComponentResponse
		ServiceError  error
		Params        component.CreateTriggerParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetComponentResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			component.NewCreateTriggerParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewCreateTriggerParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().CreateTrigger(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.CreateTrigger(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.CreateTriggerOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.CreateTriggerInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_CreateTriggerVersion() {
	cases := map[string]struct {
		ServiceReturn *dto.GetComponentVersionResponse
		ServiceError  error
		Params        component.CreateTriggerVersionParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetComponentVersionResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			component.NewCreateTriggerVersionParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewCreateTriggerVersionParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().CreateTriggerVersion(gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.CreateTriggerVersion(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.CreateTriggerVersionOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.CreateTriggerVersionInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_GetAction() {
	cases := map[string]struct {
		ServiceReturn *dto.GetComponentResponse
		ServiceError  error
		Params        component.GetActionParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetComponentResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			component.NewGetActionParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewGetActionParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetAction(gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetAction(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.GetActionOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.GetActionInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_GetActionVersion() {
	cases := map[string]struct {
		ServiceReturn *dto.GetComponentVersionResponse
		ServiceError  error
		Params        component.GetActionVersionParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetComponentVersionResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			component.NewGetActionVersionParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewGetActionVersionParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetActionVersion(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetActionVersion(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.GetActionVersionOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.GetActionVersionInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_GetAllActions() {
	cases := map[string]struct {
		ServiceReturn []*dto.GetComponentResponse
		ServiceError  error
		Params        component.GetAllActionsParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			[]*dto.GetComponentResponse{{ID: testhelper.Int64Ptr(rand.Int63())}},
			nil,
			component.NewGetAllActionsParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewGetAllActionsParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetAllActions(gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetAllActions(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.GetAllActionsOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.GetAllActionsInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_GetAllBases() {
	cases := map[string]struct {
		ServiceReturn []*dto.GetComponentResponse
		ServiceError  error
		Params        component.GetAllBasesParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			[]*dto.GetComponentResponse{{ID: testhelper.Int64Ptr(rand.Int63())}},
			nil,
			component.NewGetAllBasesParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewGetAllBasesParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetAllBases(gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetAllBases(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.GetAllBasesOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.GetAllBasesInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_GetAllTriggers() {
	cases := map[string]struct {
		ServiceReturn []*dto.GetComponentResponse
		ServiceError  error
		Params        component.GetAllTriggersParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			[]*dto.GetComponentResponse{{ID: testhelper.Int64Ptr(rand.Int63())}},
			nil,
			component.NewGetAllTriggersParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewGetAllTriggersParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetAllTriggers(gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetAllTriggers(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.GetAllTriggersOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.GetAllTriggersInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_GetBase() {
	cases := map[string]struct {
		ServiceReturn *dto.GetComponentResponse
		ServiceError  error
		Params        component.GetBaseParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetComponentResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			component.NewGetBaseParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewGetBaseParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetBase(gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetBase(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.GetBaseOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.GetBaseInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_GetBaseVersion() {
	cases := map[string]struct {
		ServiceReturn *dto.GetBaseVersionResponse
		ServiceError  error
		Params        component.GetBaseVersionParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetBaseVersionResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			component.NewGetBaseVersionParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewGetBaseVersionParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetBaseVersion(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetBaseVersion(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.GetBaseVersionOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.GetBaseVersionInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_GetTrigger() {
	cases := map[string]struct {
		ServiceReturn *dto.GetComponentResponse
		ServiceError  error
		Params        component.GetTriggerParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetComponentResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			component.NewGetTriggerParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewGetTriggerParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetTrigger(gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetTrigger(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.GetTriggerOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.GetTriggerInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_GetTriggerVersion() {
	cases := map[string]struct {
		ServiceReturn *dto.GetComponentVersionResponse
		ServiceError  error
		Params        component.GetTriggerVersionParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetComponentVersionResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			component.NewGetTriggerVersionParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewGetTriggerVersionParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetTriggerVersion(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetTriggerVersion(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.GetTriggerVersionOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.GetTriggerVersionInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_UpdateAction() {
	cases := map[string]struct {
		ServiceReturn *dto.GetComponentResponse
		ServiceError  error
		Params        component.UpdateActionParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetComponentResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			component.NewUpdateActionParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewUpdateActionParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().UpdateAction(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.UpdateAction(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.UpdateActionOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.UpdateActionInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_UpdateBase() {
	cases := map[string]struct {
		ServiceReturn *dto.GetComponentResponse
		ServiceError  error
		Params        component.UpdateBaseParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetComponentResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			component.NewUpdateBaseParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewUpdateBaseParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().UpdateBase(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.UpdateBase(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.UpdateBaseOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.UpdateBaseInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *ComponentTestSuite) TestComponentHandler_UpdateTrigger() {
	cases := map[string]struct {
		ServiceReturn *dto.GetComponentResponse
		ServiceError  error
		Params        component.UpdateTriggerParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetComponentResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			component.NewUpdateTriggerParams(),
			s.factory.DTOUser(false),
			true},
		"service returns an error": {
			nil,
			errors.New("failed"),
			component.NewUpdateTriggerParams(),
			s.factory.DTOUser(false),
			false},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().UpdateTrigger(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.UpdateTrigger(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*component.UpdateTriggerOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*component.UpdateTriggerInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

// This will make sure the test suite will run
// Don't put any logic for setting up the tests in here, use the hooks from the test suite for that
func TestComponentHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ComponentTestSuite))
}
