package pipeline_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/factory"
	"bitbucket.org/unchain/ares/pkg/testhelper"
	"github.com/stretchr/testify/suite"
)

/*
This contains the shared stuff for all pipeline service tests
*/

type TestSuite struct {
	suite.Suite
	factory *factory.Factory
	helper  *testhelper.Helper
	ares    *ares.Server
	service ares.PipelineService
	cleanup func()
}

// This runs just once before all tests
func (s *TestSuite) SetupSuite() {
	s.factory = factory.NewFactory(&s.Suite)
	s.ares, s.cleanup = s.factory.AresFactory()
	s.helper = testhelper.NewHelper(&s.Suite, s.ares)
	s.service = s.ares.PipelineService

	s.factory.SetAres(s.ares)
}

func (s *TestSuite) TearDownSuite() {
	s.cleanup()
}

// This will make sure the test suite will run
// Don't put any logic for setting up the tests in here, use the hooks from the test suite for that
func TestPipelineService(t *testing.T) {
	testhelper.SkipInBitbucket(t)

	suite.Run(t, new(TestSuite))
}
