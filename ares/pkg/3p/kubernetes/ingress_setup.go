package kubernetes

import (
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/unchainio/pkg/errors"
	"k8s.io/api/extensions/v1beta1"
	errorv1 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (s *Service) setupIngress(params *ares.DeploymentParams) (*v1beta1.Ingress, error) {
	var tls []v1beta1.IngressTLS

	if s.cfg.TLS {
		tls = []v1beta1.IngressTLS{
			{
				Hosts: []string{
					params.Host,
				},
				SecretName: params.Host,
			},
		}
	}

	o := &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:   params.FullName,
			Labels: getLabels(params),
			Annotations: map[string]string{
				"certmanager.k8s.io/cluster-issuer":          "letsencrypt-prod",
				"nginx.ingress.kubernetes.io/rewrite-target": params.RewriteTarget,
			},
		},
		Spec: v1beta1.IngressSpec{
			TLS: tls,
			Rules: []v1beta1.IngressRule{
				{
					Host: params.Host,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: params.Path,
									Backend: v1beta1.IngressBackend{
										ServiceName: params.FullName,
										ServicePort: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 80,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	c := s.Client.ExtensionsV1beta1().Ingresses(params.OrganizationName)
	res, err := c.Update(o)

	switch {
	case err == nil:
		s.log.Debugf("updated ingress %q.\n", res.GetObjectMeta().GetName())
	case errorv1.IsNotFound(err):
		res, err = c.Create(o)
		if err != nil {
			return res, errors.Wrap(err, "could not create ingress")
		}
		s.log.Debugf("ingress created")
	default:
		return res, errors.Wrap(err, "could not update nor create ingress")
	}

	return res, nil
}
