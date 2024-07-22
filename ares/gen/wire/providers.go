package wire

import (
	"testing"

	mock_ares "bitbucket.org/unchain/ares/gen/mocks"
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/golang/mock/gomock"
)

func MockAWSProvider(t *testing.T) ares.AWSClient {
	return mock_ares.NewMockAWSClient(gomock.NewController(t))
}

func NilAWSProvider() ares.AWSClient {
	return nil
}

func MockSubscriptionProvider(t *testing.T) ares.SubscriptionService {
	return mock_ares.NewMockSubscriptionService(gomock.NewController(t))
}

func MockKubernetesProvider() *mock_ares.MockDeploymentService {
	return new(mock_ares.MockDeploymentService)
}
