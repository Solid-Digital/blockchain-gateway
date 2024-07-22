package stub

import (
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchainio/interfaces/logger"
)

type Stub struct {
	logger.Logger
}

func New(logger logger.Logger) domain.Stub {
	return &Stub{logger}
}
