package event_test

import (
	"testing"

	"bitbucket.org/unchain/ethereum2/pkg/helper"

	"bitbucket.org/unchain/ethereum2/pkg/factory"
	"github.com/stretchr/testify/suite"
	"github.com/unchainio/interfaces/logger"
)

type TestSuite struct {
	suite.Suite
	logger  logger.Logger
	helper  *helper.Helper
	factory *factory.Factory
}

func (s *TestSuite) SetupSuite() {
	s.logger = factory.DefaultLogger(&s.Suite)
	s.helper = helper.NewHelper(&s.Suite, s.logger)
	s.factory = factory.NewFactory(&s.Suite, s.logger, s.helper)
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
