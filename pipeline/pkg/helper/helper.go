package helper

import (
	"github.com/stretchr/testify/suite"
	"github.com/unchainio/interfaces/logger"
)

type Helper struct {
	suite  *suite.Suite
	logger logger.Logger
}

func NewHelper(s *suite.Suite, logger logger.Logger) *Helper {
	return &Helper{suite: s, logger: logger}
}
