package action_test

import (
	"testing"

	"bitbucket.org/unchain/ethereum2/pkg/factory"
	"bitbucket.org/unchain/ethereum2/pkg/helper"
	"github.com/stretchr/testify/suite"
	"github.com/unchainio/interfaces/logger"
)

const (
	DefaultAccount       = "0xfdfa8d41f986c80904bf4825402e788f3121e7af"
	NonRegisteredAccount = "0xfac399e49f5b6867af186390270af252e683b154"

	NonExistingContractAddress = "0xd0a6e6c54dbc68db5db3a091b171a77407ff7ccf"
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
