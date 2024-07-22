package cron_trigger

import (
	"github.com/robfig/cron/v3"
	"github.com/unchain/pipeline/pkg/domain"
)

type Trigger struct {
	config         *Config
	stub           domain.Stub
	cron           *cron.Cron
	RequestChannel chan *domain.Request
}

func NewTrigger() *Trigger {
	return &Trigger{}
}
