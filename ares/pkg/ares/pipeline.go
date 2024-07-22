package ares

import (
	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/pkg/3p/apperr"
)

type PipelineService interface {
	CreatePipeline(params *dto.CreatePipelineRequest, orgName string, principal *dto.User) (*dto.GetPipelineResponse, *apperr.Error)
	UpdatePipeline(params *dto.UpdatePipelineRequest, orgName string, pipelineName string, principal *dto.User) (*dto.GetPipelineResponse, *apperr.Error)
	GetPipeline(orgName string, pipelineName string) (*dto.GetPipelineResponse, *apperr.Error)
	DeletePipeline(orgName string, pipelineName string, principal *dto.User) *apperr.Error
	UpdateDraftConfiguration(params *dto.UpdateDraftConfigurationRequest, orgName string, pipelineName string, principal *dto.User) (*dto.GetConfigurationResponse, *apperr.Error)
	SaveDraftConfigurationAsNew(params *dto.SaveDraftConfigurationAsNewRequest, orgName string, pipelineName string, principal *dto.User) (*dto.GetConfigurationResponse, *apperr.Error)
	GetAllPipelines(orgName string) (dto.GetAllPipelinesResponse, *apperr.Error)
	GetConfiguration(orgName string, pipelineName string, revision int64) (*dto.GetConfigurationResponse, *apperr.Error)
	GetDraftConfiguration(orgName string, pipelineName string) (*dto.GetConfigurationResponse, *apperr.Error)
	DeployConfiguration(params *dto.DeployConfigurationRequest, orgName string, pipelineName string, envName string, user *dto.User) (*dto.GetDeploymentResponse, *apperr.Error)
	RemoveDeployment(orgName string, pipelineName string, envName string) *apperr.Error
	GetDeploymentLogs(orgName string, pipelineName string, envName string, from string, to string, limit string) ([]*dto.LogLine, *apperr.Error)
	GetDeployment(orgName string, pipelineName string, envName string) (*dto.GetDeploymentResponse, *apperr.Error)
	CreateEnvironmentVariable(params *dto.CreateEnvironmentVariableRequest, orgName string, pipelineName string, envName string, user *dto.User) (*dto.GetEnvironmentVariableResponse, *apperr.Error)
	GetAllEnvironmentVariables(orgName string, pipelineName string, envName string, user *dto.User) (dto.GetAllEnvironmentVariablesResponse, *apperr.Error)
	DeleteEnvironmentVariable(orgName string, pipelineName string, envName string, varID int64, user *dto.User) *apperr.Error
	UpdateEnvironmentVariable(params *dto.UpdateEnvironmentVariablesRequest, orgName string, pipelineName string, envName string, varID int64, user *dto.User) (*dto.GetEnvironmentVariableResponse, *apperr.Error)
}
