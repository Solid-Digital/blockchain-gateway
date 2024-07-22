package pipeline

func SetFetchAdapterLogsFn(fn func(s *Service, orgName, pipelineName string, envName string, from, to string) (DoFn, ClearFn, error)) {
	fetchDeploymentLogsFn = fn
}

func ResetFetchAdapterLogsFn() {
	fetchDeploymentLogsFn = fetchDeploymentLogs
}
