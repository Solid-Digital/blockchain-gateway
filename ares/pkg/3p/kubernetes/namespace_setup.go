package kubernetes

import (
	"github.com/unchainio/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	errorv1 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *Service) setupNamespace(name string) (*apiv1.Namespace, error) {
	res, err := s.Client.CoreV1().Namespaces().Create(
		&apiv1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"field.cattle.io/projectId": s.cfg.RancherProjectID,
				},
				Annotations: map[string]string{
					"field.cattle.io/projectId": s.cfg.RancherClusterID + ":" + s.cfg.RancherProjectID,
				},
				Name: name,
			},
		},
	)
	if err != nil && !errorv1.IsAlreadyExists(err) {
		return nil, errors.Wrap(err, "")
	}

	return res, nil
}
