package kubernetes_test

import (
	v1 "k8s.io/api/apps/v1"
)

func (s *TestSuite) TestService_GetAllDeployments() {
	org, _, _, deployment1 := s.factory.DeploymentFromService()
	_, _, deployment2 := s.factory.DeploymentFromServiceForOrg(org)
	_, _, deployment3 := s.factory.DeploymentFromServiceForOrg(org)

	list, err := s.ares.DeploymentService.GetAllDeployments(org.Name, []string{deployment1.FullName, deployment2.FullName})

	s.Require().NoError(err)
	s.Require().Equal(2, len(list.Items))
	s.Require().True(contains(list.Items, deployment1.FullName))
	s.Require().True(contains(list.Items, deployment2.FullName))
	s.Require().False(contains(list.Items, deployment3.FullName))
}

func contains(items []v1.Deployment, fullName string) bool {
	for _, item := range items {
		if item.Name == fullName {
			return true
		}
	}

	return false
}
