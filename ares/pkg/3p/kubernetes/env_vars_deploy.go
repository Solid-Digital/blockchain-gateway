package kubernetes

import (
	"bitbucket.org/unchain/ares/pkg/ares"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//
func (s *Service) deployEnvVars(params *ares.DeploymentParams) (vars *apiv1.Secret, err error) {
	data := make(map[string]string)

	for _, v := range params.EnvironmentVariables {
		data[v.Key] = v.Value
	}

	vars = &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:   params.FullName,
			Labels: getLabels(params),
		},
		StringData: data,
	}

	s.log.Debugf("deploying secret containing environment variables...")
	err = s.upsertSecret(params.OrganizationName, vars)
	if err != nil {
		return nil, err
	}

	return vars, nil
}
