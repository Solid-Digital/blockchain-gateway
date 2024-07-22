package fileparser_action_test

import (
	"github.com/stretchr/testify/suite"
	"github.com/unchain/pipeline/pkg/factory"
	"github.com/unchain/pipeline/pkg/helper"
	"github.com/unchainio/interfaces/logger"
	"testing"
)

type TestSuite struct {
	suite.Suite
	logger  logger.Logger
	helper  *helper.Helper
}

func (s *TestSuite) SetupSuite() {
	s.logger = factory.DefaultLogger(&s.Suite)
	s.helper = helper.NewHelper(&s.Suite, s.logger)
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}