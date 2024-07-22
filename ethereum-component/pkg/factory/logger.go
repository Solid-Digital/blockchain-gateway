package factory

import (
	"github.com/stretchr/testify/suite"
	"github.com/unchainio/interfaces/logger"
	"github.com/unchainio/pkg/xlogger"
)

func DefaultLogger(s *suite.Suite) logger.Logger {
	cfg := new(xlogger.Config)
	logger, err := xlogger.New(cfg)

	s.Require().NoError(err)

	return logger
}
