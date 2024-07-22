package kubernetes

import (
	"fmt"
	"time"

	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/unchainio/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	errorv1 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (s *Service) CreateDeployment(params *ares.DeploymentParams) (*appsv1.Deployment, error) {
	var err error

	_, err = s.setupNamespace(params.OrganizationName)
	if err != nil {
		return nil, err
	}

	_, err = s.setupConfigMap(params)
	if err != nil {
		return nil, err
	}

	_, err = s.deployEnvVars(params)
	if err != nil {
		return nil, err
	}

	deployment, err := s.setupDeployment(params)
	if err != nil {
		return nil, err
	}

	_, err = s.setupService(params)
	if err != nil {
		return nil, err
	}

	_, err = s.setupIngress(params)
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func (s *Service) setupDeployment(params *ares.DeploymentParams) (*appsv1.Deployment, error) {
	s.log.Debugf("Deploying `%s` to kubernetes", params.Image)

	var regcred []apiv1.LocalObjectReference

	if s.cfg.RegistryCredentialsSecret != "" {
		regcred = []apiv1.LocalObjectReference{{
			Name: s.cfg.RegistryCredentialsSecret,
		}}
	}

	o := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        params.FullName,
			Labels:      getLabels(params),
			Annotations: getAnnotations(),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"fullName": params.FullName},
			},
			Replicas: int32Ptr(int(params.Replicas)),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      getLabels(params),
					Annotations: getAnnotations(),
				},

				Spec: apiv1.PodSpec{
					Volumes: []apiv1.Volume{
						{
							Name: params.FullName,
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: params.FullName,
									},
								},
							},
						},
					},
					Containers: []apiv1.Container{
						{
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      params.FullName,
									ReadOnly:  true,
									MountPath: "/etc/opt/unchain/",
								},
							},
							Name:  "adapter",
							Image: params.Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
							EnvFrom: []apiv1.EnvFromSource{
								{
									SecretRef: &apiv1.SecretEnvSource{
										LocalObjectReference: apiv1.LocalObjectReference{
											Name: params.FullName,
										},
									},
								},
							},
							ImagePullPolicy: apiv1.PullAlways,
						},
					},
					ImagePullSecrets: regcred,
					RestartPolicy:    apiv1.RestartPolicyAlways,
					DNSPolicy:        apiv1.DNSClusterFirst,
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(0),
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(1),
					},
				},
			},
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1",
		},
	}

	s.log.Debugf("creating deployment...")

	// Implement deployment update-or-create semantics.
	c := s.Client.AppsV1().Deployments(params.OrganizationName)
	res, err := c.Update(o)
	switch {
	case err == nil:
		s.log.Debugf("updated deployment %q.\n", res.GetObjectMeta().GetName())
	case !errorv1.IsNotFound(err):
		return res, errors.Wrap(err, "could not update deployment controller: ")
	default:
		res, err = c.Create(o)
		if err != nil {
			return res, errors.Wrap(err, "could not create deployment controller: ")
		}
		s.log.Debugf("deployment controller created")
	}

	return res, nil
}

func getAnnotations() map[string]string {
	return map[string]string{
		// This makes it so  that every time we update the deployment, all pods are forced to restart.
		"updated-at": time.Now().String(),
	}
}

func getLabels(params *ares.DeploymentParams) map[string]string {
	return map[string]string{
		"fullName":     params.FullName,
		"organization": params.OrganizationName,
		"environment":  params.EnvironmentName,
		"revision":     fmt.Sprintf("%d", params.Revision),
	}
}

func int64Ptr(i int64) *int64 {
	return &i
}

func int32Ptr(i int) *int32 {
	ii := int32(i)
	return &ii
}
