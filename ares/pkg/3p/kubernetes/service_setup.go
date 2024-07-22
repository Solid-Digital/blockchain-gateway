package kubernetes

import (
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/unchainio/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	errorv1 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (s *Service) setupService(params *ares.DeploymentParams) (*apiv1.Service, error) {
	o := &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   params.FullName,
			Labels: getLabels(params),
		},
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeNodePort,
			// TODO should this be hard coded?
			Selector: map[string]string{"fullName": params.FullName},
			Ports: []apiv1.ServicePort{
				{
					Protocol: apiv1.ProtocolTCP,
					Port:     80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 80,
					},
				},
			},
		},
	}

	// Implement service update-or-create methods.
	c := s.Client.CoreV1().Services(params.OrganizationName)
	res, err := c.Get(params.FullName, metav1.GetOptions{})
	switch {
	case err == nil:
		o.ObjectMeta.ResourceVersion = res.ObjectMeta.ResourceVersion
		o.Spec.ClusterIP = res.Spec.ClusterIP
		res, err = c.Update(o)
		if err != nil {
			return res, errors.Wrap(err, "failed to update service: ")
		}
		s.log.Debugf("service updated")
	case errorv1.IsNotFound(err):
		res, err = c.Create(o)
		if err != nil {
			return res, errors.Wrap(err, "failed to create service: ")
		}
		s.log.Debugf("service created and exposed")
	default:
		return res, errors.Wrap(err, "unexpected error: ")
	}
	return res, nil
}
