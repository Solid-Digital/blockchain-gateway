package ares

import (
	"bitbucket.org/unchain/ares/gen/orm"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
)

type DeploymentService interface {
	GetDeploymentPods(namespace, name string) (*apiv1.PodList, error)
	GetDeployment(namespace, name string) (*appsv1.Deployment, error)
	GetAllDeployments(namespace string, names []string) (*appsv1.DeploymentList, error)
	CreateDeployment(params *DeploymentParams) (*appsv1.Deployment, error)
	DeleteDeployment(namespace string, name string) error
	DeleteAllDeployments(namespace string, names []string) error
}

type DeploymentParams struct {
	*orm.DeploymentDTO
	Config               string
	OrganizationName     string
	EnvironmentName      string
	Revision             int64
	EnvironmentVariables []*orm.EnvironmentVariableDTO
}
