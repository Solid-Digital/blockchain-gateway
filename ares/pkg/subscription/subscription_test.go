package subscription_test

import (
	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/factory"
	"bitbucket.org/unchain/ares/pkg/testhelper"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TestSuite struct {
	suite.Suite
	factory *factory.Factory
	helper  *testhelper.Helper
	ares    *ares.Server
	service ares.SubscriptionService
	cleanup func()
}

func (s *TestSuite) SetupSuite() {
	s.factory = factory.NewFactory(&s.Suite)
	s.ares, s.cleanup = s.factory.AresFactory()
	s.helper = testhelper.NewHelper(&s.Suite, s.ares)

	s.service = s.ares.SubscriptionService

	s.factory.SetAres(s.ares)
}

func (s *TestSuite) TearDownSuite() {
	s.cleanup()
}

func TestSubscriptionService(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
