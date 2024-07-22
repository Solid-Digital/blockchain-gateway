package organization_test

import (
	"testing"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/testhelper"

	"bitbucket.org/unchain/ares/pkg/factory"
	"github.com/stretchr/testify/suite"
)

/*
This contains the shared stuff for all organization service tests
*/

type OrganizationTestSuite struct {
	suite.Suite
	factory *factory.Factory
	helper  *testhelper.Helper
	ares    *ares.Server
	service ares.OrganizationService
	cleanup func()
}

// This runs just once before all tests
func (s *OrganizationTestSuite) SetupSuite() {
	s.factory = factory.NewFactory(&s.Suite)
	s.ares, s.cleanup = s.factory.AresFactory()
	s.helper = testhelper.NewHelper(&s.Suite, s.ares)
	s.service = s.ares.OrganizationService

	s.factory.SetAres(s.ares)
}

func (s *OrganizationTestSuite) TearDownSuite() {
	s.cleanup()
}

// This will make sure the test suite will run
// Don't put any logic for setting up the tests in here, use the hooks from the test suite for that
func TestOrganizationTestSuite(t *testing.T) {
	suite.Run(t, new(OrganizationTestSuite))
}
