package trigger_test

import (
	"testing"

	"bitbucket.org/unchain/ethereum2/pkg/factory"
	"bitbucket.org/unchain/ethereum2/pkg/helper"
	"github.com/stretchr/testify/suite"
	"github.com/unchainio/interfaces/logger"
)

const (
	DefaultAccount     = "0xfdfa8d41f986c80904bf4825402e788f3121e7af"
	AlternativeAccount = "0x69091d42c8307d9a24a47b2d92d4506604ae44b9"
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
