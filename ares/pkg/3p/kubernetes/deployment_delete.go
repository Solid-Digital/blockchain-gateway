package kubernetes

import (
	"github.com/unchainio/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func (s *Service) DeleteDeployment(namespace string, name string) error {
	deletePropagation := metav1.DeletePropagationForeground
	err := s.Client.AppsV1().Deployments(namespace).Delete(name, &metav1.DeleteOptions{
		GracePeriodSeconds: int64Ptr(0),
		PropagationPolicy:  &deletePropagation,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to remove deployment `%s`", name)
	}

	err = s.Client.CoreV1().Services(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to remove service `%s`", name)
	}

	err = s.Client.ExtensionsV1beta1().Ingresses(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to remove ingress `%s`", name)
	}

	err = s.Client.CoreV1().Secrets(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to remove secrets `%s`", name)
	}

	err = s.Client.CoreV1().ConfigMaps(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to remove config map `%s`", name)
	}

	return nil
}

func (s *Service) DeleteAllDeployments(namespace string, names []string) error {
	deletePropagation := metav1.DeletePropagationForeground
	deleteOpts := &metav1.DeleteOptions{
		GracePeriodSeconds: int64Ptr(0),
		PropagationPolicy:  &deletePropagation,
	}

	req, err := labels.NewRequirement("fullName", selection.In, names)
	if err != nil {
		return errors.Wrapf(err, "failed to construct label requirement")
	}

	listOpts := &metav1.ListOptions{
		LabelSelector: labels.NewSelector().Add(*req).String(),
	}

	err = s.Client.AppsV1().Deployments(namespace).DeleteCollection(deleteOpts, *listOpts)
	if err != nil {
		return errors.Wrapf(err, "failed to remove deployments `%s`", names)
	}

	// Deleting a collection of services is not implemented
	// https://github.com/kubernetes/kubernetes/issues/68468#issuecomment-419981870
	for _, name := range names {
		err = s.Client.CoreV1().Services(namespace).Delete(name, &metav1.DeleteOptions{})
		if err != nil {
			return errors.Wrapf(err, "failed to remove service `%s`", name)
		}
	}

	err = s.Client.ExtensionsV1beta1().Ingresses(namespace).DeleteCollection(&metav1.DeleteOptions{}, *listOpts)
	if err != nil {
		return errors.Wrapf(err, "failed to remove ingresses `%s`", names)
	}

	err = s.Client.CoreV1().Secrets(namespace).DeleteCollection(&metav1.DeleteOptions{}, *listOpts)
	if err != nil {
		return errors.Wrapf(err, "failed to remove secrets `%s`", names)
	}

	err = s.Client.CoreV1().ConfigMaps(namespace).DeleteCollection(&metav1.DeleteOptions{}, *listOpts)
	if err != nil {
		return errors.Wrapf(err, "failed to remove config maps `%s`", names)
	}

	return nil
}
