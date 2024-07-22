package api_trigger_test

import (
	"github.com/unchain/pipeline/pkg/factory"
	"github.com/unchain/pipeline/pkg/helper"
	"github.com/unchain/pipeline/pkg/triggers/api_trigger/test_helpers"
	"testing"

	"github.com/unchainio/interfaces/logger"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	logger       logger.Logger
	helper       *helper.Helper
	factory      *factory.Factory
	test_helpers *test_helpers.TestHelpers
}

func (s *TestSuite) SetupSuite() {
	s.logger = factory.DefaultLogger(&s.Suite)
	s.helper = helper.NewHelper(&s.Suite, s.logger)
	s.factory = factory.NewFactory(&s.Suite, s.logger, s.helper)
	s.test_helpers = test_helpers.NewTestHelpers(&s.Suite, s.logger, s.helper)
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
