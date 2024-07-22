package kubernetes

import (
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/unchainio/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	errorv1 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *Service) setupConfigMap(params *ares.DeploymentParams) (*apiv1.ConfigMap, error) {
	o := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:   params.FullName,
			Labels: getLabels(params),
		},
		Data: map[string]string{
			"config.toml": params.Config,
		},
	}

	s.log.Debugf("creating config map...")

	// Implement config map update-or-create semantics.
	c := s.Client.CoreV1().ConfigMaps(params.OrganizationName)
	res, err := c.Update(o)
	switch {
	case err == nil:
		s.log.Debugf("updated config map %q.\n", res.GetObjectMeta().GetName())
	case !errorv1.IsNotFound(err):
		return res, errors.Wrap(err, "could not update config map: ")
	default:
		res, err = c.Create(o)
		if err != nil {
			return res, errors.Wrap(err, "could not create deployment controller: ")
		}
		s.log.Debugf("config map created")
	}

	return res, nil
}
