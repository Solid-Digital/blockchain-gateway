package kubernetes_test

import (
	"os"
	"testing"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/factory"
	"bitbucket.org/unchain/ares/pkg/testhelper"
	"github.com/stretchr/testify/suite"
)

/*
This contains the shared stuff for all kube service tests
*/

type TestSuite struct {
	suite.Suite
	factory *factory.Factory
	helper  *testhelper.Helper
	ares    *ares.Server
	cleanup func()
}

// This runs just once before all tests
func (s *TestSuite) SetupSuite() {
	// FIXME(e-nikolov) do this so that henk's config loader works. Currently it assumes that the tests are run from 2 levels above the config.toml
	os.Chdir("../")
	s.factory = factory.NewFactory(&s.Suite)
	s.ares, s.cleanup = s.factory.AresFactory()
	s.helper = testhelper.NewHelper(&s.Suite, s.ares)

	s.factory.SetAres(s.ares)
}

func (s *TestSuite) TearDownSuite() {
	s.cleanup()
}

// This will make sure the test suite will run
// Don't put any logic for setting up the tests in here, use the hooks from the test suite for that
func TestTestSuite(t *testing.T) {
	testhelper.SkipInBitbucket(t)

	suite.Run(t, new(TestSuite))
}
