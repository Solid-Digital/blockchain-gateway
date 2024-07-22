package janus

import (
	"bitbucket.org/unchain/ares/pkg/pipeline/internal/janus/pkg/pipeline"
	"github.com/unchainio/pkg/xlogger"
)

type Config struct {
	Logger   *xlogger.Config
	Pipeline *pipeline.Config
}
