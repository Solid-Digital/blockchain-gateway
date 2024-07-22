package http_test

import (
	"math/rand"
	"testing"

	"bitbucket.org/unchain/ares/pkg/testhelper"

	"bitbucket.org/unchain/ares/gen/api/operations/pipeline"
	"bitbucket.org/unchain/ares/gen/dto"
	"github.com/stretchr/testify/require"
	"github.com/unchainio/pkg/errors"

	mock_ares "bitbucket.org/unchain/ares/gen/mocks"
	"bitbucket.org/unchain/ares/pkg/factory"
	"bitbucket.org/unchain/ares/pkg/http"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type PipelineTestSuite struct {
	suite.Suite
	factory       *factory.Factory
	handler       *http.PipelineHandler
	mockedService *mock_ares.MockPipelineService
}

func (s *PipelineTestSuite) SetupSuite() {
	s.T().Skip()
	s.factory = factory.NewFactory(&s.Suite)
}

// This runs before each test
// A new mock needs to be created for each test, otherwise it will fail
func (s *PipelineTestSuite) SetupTest() {
	service := mock_ares.NewMockPipelineService(gomock.NewController(s.T()))
	s.mockedService = service
	s.handler = http.NewPipelineHandler(service, s.factory.Logger())
}

func (s *PipelineTestSuite) TestPipelineHandler_CreatePipeline() {
	cases := map[string]struct {
		ServiceReturn *dto.GetPipelineResponse
		ServiceError  error
		Params        pipeline.CreatePipelineParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetPipelineResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			pipeline.NewCreatePipelineParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			nil,
			errors.New("failed"),
			pipeline.NewCreatePipelineParams(),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().CreatePipeline(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.CreatePipeline(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*pipeline.CreatePipelineOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*pipeline.CreatePipelineInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *PipelineTestSuite) TestPipelineHandler_DeletePipeline() {
	cases := map[string]struct {
		ServiceError error
		Params       pipeline.DeletePipelineParams
		Principal    *dto.User
		Success      bool
	}{
		"service returns no error": {
			nil,
			pipeline.NewDeletePipelineParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			errors.New("failed"),
			pipeline.NewDeletePipelineParams(),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().DeletePipeline(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceError)
			response := s.handler.DeletePipeline(tc.Params, tc.Principal)

			if tc.Success {
				_, ok := response.(*pipeline.DeletePipelineNoContent)
				require.True(t, ok)
			} else {
				result, ok := response.(*pipeline.DeletePipelineInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *PipelineTestSuite) TestPipelineHandler_DeployConfiguration() {
	cases := map[string]struct {
		ServiceReturn *dto.GetDeploymentResponse
		ServiceError  error
		Params        pipeline.DeployConfigurationParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetDeploymentResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			pipeline.NewDeployConfigurationParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			nil,
			errors.New("failed"),
			pipeline.NewDeployConfigurationParams(),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().DeployConfiguration(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.DeployConfiguration(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*pipeline.DeployConfigurationOK)
				require.Equal(t, tc.ServiceReturn, result.Payload)
				require.True(t, ok)
			} else {
				result, ok := response.(*pipeline.DeployConfigurationInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *PipelineTestSuite) TestPipelineHandler_GetAllPipelines() {
	cases := map[string]struct {
		ServiceReturn dto.GetAllPipelinesResponse
		ServiceError  error
		Params        pipeline.GetAllPipelinesParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			dto.GetAllPipelinesResponse{{ID: testhelper.Int64Ptr(rand.Int63())}},
			nil,
			pipeline.NewGetAllPipelinesParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			nil,
			errors.New("failed"),
			pipeline.NewGetAllPipelinesParams(),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetAllPipelines(gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetAllPipelines(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*pipeline.GetAllPipelinesOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*pipeline.GetAllPipelinesInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *PipelineTestSuite) TestPipelineHandler_GetConfiguration() {
	cases := map[string]struct {
		ServiceReturn *dto.GetConfigurationResponse
		ServiceError  error
		Params        pipeline.GetConfigurationParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetConfigurationResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			pipeline.NewGetConfigurationParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			nil,
			errors.New("failed"),
			pipeline.NewGetConfigurationParams(),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetConfiguration(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetConfiguration(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*pipeline.GetConfigurationOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*pipeline.GetConfigurationInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *PipelineTestSuite) TestPipelineHandler_GetDeployment() {
	cases := map[string]struct {
		ServiceReturn *dto.GetDeploymentResponse
		ServiceError  error
		Params        pipeline.GetDeploymentParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetDeploymentResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			pipeline.NewGetDeploymentParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			nil,
			errors.New("failed"),
			pipeline.NewGetDeploymentParams(),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetDeployment(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetDeployment(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*pipeline.GetDeploymentOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*pipeline.GetDeploymentInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *PipelineTestSuite) TestPipelineHandler_GetDeploymentLogs() {
	cases := map[string]struct {
		ServiceReturn []*dto.LogLine
		ServiceError  error
		Params        pipeline.GetDeploymentLogsParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			[]*dto.LogLine{{Caller: string(rand.Int())}},
			nil,
			pipeline.NewGetDeploymentLogsParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			nil,
			errors.New("failed"),
			pipeline.NewGetDeploymentLogsParams(),
			s.factory.DTOUser(false),
			false,
		},
		"with non nil parameters": {
			[]*dto.LogLine{{Caller: string(rand.Int())}},
			nil,
			pipeline.GetDeploymentLogsParams{
				From:  testhelper.StringPtr("1"),
				To:    testhelper.StringPtr("2"),
				Limit: testhelper.StringPtr("3"),
			},
			s.factory.DTOUser(false),
			true,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetDeploymentLogs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetDeploymentLogs(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*pipeline.GetDeploymentLogsOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*pipeline.GetDeploymentLogsInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *PipelineTestSuite) TestPipelineHandler_GetDraftConfiguration() {
	cases := map[string]struct {
		ServiceReturn *dto.GetConfigurationResponse
		ServiceError  error
		Params        pipeline.GetDraftConfigurationParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetConfigurationResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			pipeline.NewGetDraftConfigurationParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			nil,
			errors.New("failed"),
			pipeline.NewGetDraftConfigurationParams(),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetDraftConfiguration(gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetDraftConfiguration(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*pipeline.GetDraftConfigurationOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*pipeline.GetDraftConfigurationInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *PipelineTestSuite) TestPipelineHandler_GetPipeline() {
	cases := map[string]struct {
		ServiceReturn *dto.GetPipelineResponse
		ServiceError  error
		Params        pipeline.GetPipelineParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetPipelineResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			pipeline.NewGetPipelineParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			nil,
			errors.New("failed"),
			pipeline.NewGetPipelineParams(),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetPipeline(gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetPipeline(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*pipeline.GetPipelineOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*pipeline.GetPipelineInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *PipelineTestSuite) TestPipelineHandler_CreateEnvironmentVariable() {
	cases := map[string]struct {
		ServiceReturn *dto.GetEnvironmentVariableResponse
		ServiceError  error
		Params        pipeline.CreateEnvironmentVariableParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetEnvironmentVariableResponse{ID: rand.Int63()},
			nil,
			pipeline.NewCreateEnvironmentVariableParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			nil,
			errors.New("failed"),
			pipeline.NewCreateEnvironmentVariableParams(),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().CreateEnvironmentVariable(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.CreateEnvironmentVariable(tc.Params, tc.Principal)

			if tc.Success {
				_, ok := response.(*pipeline.CreateEnvironmentVariableOK)
				require.True(t, ok)
			} else {
				result, ok := response.(*pipeline.CreateEnvironmentVariableInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *PipelineTestSuite) TestPipelineHandler_GetAllEnvironmentVariables() {
	cases := map[string]struct {
		ServiceReturn dto.GetAllEnvironmentVariablesResponse
		ServiceError  error
		Params        pipeline.GetAllEnvironmentVariablesParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			dto.GetAllEnvironmentVariablesResponse{{ID: rand.Int63()}},
			nil,
			pipeline.NewGetAllEnvironmentVariablesParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			nil,
			errors.New("failed"),
			pipeline.NewGetAllEnvironmentVariablesParams(),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().GetAllEnvironmentVariables(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.GetAllEnvironmentVariables(tc.Params, tc.Principal)

			if tc.Success {
				_, ok := response.(*pipeline.GetAllEnvironmentVariablesOK)
				require.True(t, ok)
			} else {
				result, ok := response.(*pipeline.GetAllEnvironmentVariablesInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *PipelineTestSuite) TestPipelineHandler_DeleteEnvironmentVariable() {
	cases := map[string]struct {
		ServiceError error
		Params       pipeline.DeleteEnvironmentVariableParams
		Principal    *dto.User
		Success      bool
	}{
		"service returns no error": {
			nil,
			pipeline.NewDeleteEnvironmentVariableParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			errors.New("failed"),
			pipeline.NewDeleteEnvironmentVariableParams(),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().DeleteEnvironmentVariable(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceError)
			response := s.handler.DeleteEnvironmentVariable(tc.Params, tc.Principal)

			if tc.Success {
				_, ok := response.(*pipeline.DeleteEnvironmentVariableNoContent)
				require.True(t, ok)
			} else {
				result, ok := response.(*pipeline.DeleteEnvironmentVariableInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *PipelineTestSuite) TestPipelineHandler_UpdateEnvironmentVariable() {
	cases := map[string]struct {
		ServiceReturn *dto.GetEnvironmentVariableResponse
		ServiceError  error
		Params        pipeline.UpdateEnvironmentVariableParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetEnvironmentVariableResponse{ID: rand.Int63()},
			nil,
			pipeline.NewUpdateEnvironmentVariableParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			nil,
			errors.New("failed"),
			pipeline.NewUpdateEnvironmentVariableParams(),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().UpdateEnvironmentVariable(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.UpdateEnvironmentVariable(tc.Params, tc.Principal)

			if tc.Success {
				_, ok := response.(*pipeline.UpdateEnvironmentVariableOK)
				require.True(t, ok)
			} else {
				result, ok := response.(*pipeline.UpdateEnvironmentVariableInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *PipelineTestSuite) TestPipelineHandler_RemoveDeployment() {
	cases := map[string]struct {
		ServiceError error
		Params       pipeline.RemoveDeploymentParams
		Principal    *dto.User
		Success      bool
	}{
		"service returns no error": {
			nil,
			pipeline.NewRemoveDeploymentParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			errors.New("failed"),
			pipeline.NewRemoveDeploymentParams(),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().RemoveDeployment(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceError)
			response := s.handler.RemoveDeployment(tc.Params, tc.Principal)

			if tc.Success {
				_, ok := response.(*pipeline.RemoveDeploymentNoContent)
				require.True(t, ok)
			} else {
				result, ok := response.(*pipeline.RemoveDeploymentInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *PipelineTestSuite) TestPipelineHandler_SaveDraftConfigurationAsNew() {
	cases := map[string]struct {
		ServiceReturn *dto.GetConfigurationResponse
		ServiceError  error
		Params        pipeline.SaveDraftConfigurationAsNewParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetConfigurationResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			pipeline.NewSaveDraftConfigurationAsNewParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			nil,
			errors.New("failed"),
			pipeline.NewSaveDraftConfigurationAsNewParams(),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().SaveDraftConfigurationAsNew(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.SaveDraftConfigurationAsNew(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*pipeline.SaveDraftConfigurationAsNewOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*pipeline.SaveDraftConfigurationAsNewInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *PipelineTestSuite) TestPipelineHandler_UpdateDraftConfiguration() {
	cases := map[string]struct {
		ServiceReturn *dto.GetConfigurationResponse
		ServiceError  error
		Params        pipeline.UpdateDraftConfigurationParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetConfigurationResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			pipeline.NewUpdateDraftConfigurationParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			nil,
			errors.New("failed"),
			pipeline.NewUpdateDraftConfigurationParams(),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().UpdateDraftConfiguration(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.UpdateDraftConfiguration(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*pipeline.UpdateDraftConfigurationOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*pipeline.UpdateDraftConfigurationInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

func (s *PipelineTestSuite) TestPipelineHandler_UpdatePipeline() {
	cases := map[string]struct {
		ServiceReturn *dto.GetPipelineResponse
		ServiceError  error
		Params        pipeline.UpdatePipelineParams
		Principal     *dto.User
		Success       bool
	}{
		"service returns no error": {
			&dto.GetPipelineResponse{ID: testhelper.Int64Ptr(rand.Int63())},
			nil,
			pipeline.NewUpdatePipelineParams(),
			s.factory.DTOUser(false),
			true,
		},
		"service returns an error": {
			nil,
			errors.New("failed"),
			pipeline.NewUpdatePipelineParams(),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.SetupTest()

		s.T().Run(name, func(t *testing.T) {
			s.mockedService.EXPECT().UpdatePipeline(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.ServiceReturn, tc.ServiceError)
			response := s.handler.UpdatePipeline(tc.Params, tc.Principal)

			if tc.Success {
				result, ok := response.(*pipeline.UpdatePipelineOK)
				require.True(t, ok)
				require.Equal(t, tc.ServiceReturn, result.Payload)
			} else {
				result, ok := response.(*pipeline.UpdatePipelineInternalServerError)
				require.True(t, ok)
				require.Contains(t, result.Payload, "failed")
			}
		})
	}
}

// This will make sure the test suite will run
// Don't put any logic for setting up the tests in here, use the hooks from the test suite for that
func TestPipelineTestSuite(t *testing.T) {
	suite.Run(t, new(PipelineTestSuite))
}
