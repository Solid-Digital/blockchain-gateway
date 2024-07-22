package factory

import (
	"github.com/unchainio/interfaces/logger"
	"github.com/unchainio/pkg/xlogger"
)

func (f *Factory) Logger() logger.Logger {
	//cfg := new(xlogger.Config)
	//log, err := xlogger.New(cfg)

	//f.suite.Require().NoError(err)

	return &xlogger.Mock{}
}
