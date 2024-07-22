package kubernetes

import (
	"github.com/unchainio/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func (s *Service) GetDeployment(namespace, name string) (*appsv1.Deployment, error) {
	var err error
	c := s.Client.AppsV1().Deployments(namespace)
	deployment, err := c.Get(name, metav1.GetOptions{})

	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return deployment, nil
}

func (s *Service) GetDeploymentPods(namespace, name string) (*apiv1.PodList, error) {
	var err error

	deployment, err := s.GetDeployment(namespace, name)
	if err != nil {
		return nil, err
	}

	c := s.Client.CoreV1().Pods(namespace)
	set := labels.Set(deployment.Spec.Selector.MatchLabels)

	pod, err := c.List(metav1.ListOptions{LabelSelector: set.AsSelector().String()})

	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return pod, nil
}

func (s *Service) GetAllDeployments(namespace string, names []string) (*appsv1.DeploymentList, error) {
	req, err := labels.NewRequirement("fullName", selection.In, names)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to construct label requirement")
	}

	listOpts := metav1.ListOptions{
		LabelSelector: labels.NewSelector().Add(*req).String(),
	}

	list, err := s.Client.AppsV1().Deployments(namespace).List(listOpts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get deployment list")
	}

	return list, err
}
