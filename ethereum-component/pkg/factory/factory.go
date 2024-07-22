package factory

import (
	"bitbucket.org/unchain/ethereum2/pkg/helper"
	"github.com/stretchr/testify/suite"
	"github.com/unchainio/interfaces/logger"
)

type Factory struct {
	suite  *suite.Suite
	logger logger.Logger
	helper *helper.Helper
}

func NewFactory(s *suite.Suite, logger logger.Logger, helper *helper.Helper) *Factory {
	return &Factory{suite: s, logger: logger, helper: helper}
}
