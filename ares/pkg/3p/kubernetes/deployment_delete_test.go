package kubernetes_test

func (s *TestSuite) TestService_DeleteAllDeployments() {
	org, pipeline1, _, deployment1 := s.factory.DeploymentFromService()
	pipeline2, _, deployment2 := s.factory.DeploymentFromServiceForOrg(org)
	pipeline3, _, deployment3 := s.factory.DeploymentFromServiceForOrg(org)

	err := s.ares.DeploymentService.DeleteAllDeployments(org.Name, []string{deployment1.FullName, deployment2.FullName})
	s.Require().NoError(err)

	// Pipeline1 should be removed
	s.Require().False(s.helper.HasActiveDeployments(pipeline1))
	s.Require().False(s.helper.HasActivePods(org.Name, deployment1.FullName))

	// Pipeline2 should be removed
	s.Require().False(s.helper.HasActiveDeployments(pipeline2))
	s.Require().False(s.helper.HasActivePods(org.Name, deployment2.FullName))

	// Pipeline3 should not be removed
	s.Require().True(s.helper.HasActiveDeployments(pipeline3))
	s.Require().True(s.helper.HasActivePods(org.Name, deployment3.FullName))
}
