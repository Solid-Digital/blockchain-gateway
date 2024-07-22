package amqp_trigger

import (
	"github.com/streadway/amqp"
	"github.com/unchain/pipeline/pkg/domain"
	"sync"
)

type Trigger struct {
	config             *Config
	stub               domain.Stub
	amqpConn           *amqp.Connection
	amqpChannel        *amqp.Channel
	RequestChannel     chan *domain.Request
	ResponseChannelMap *sync.Map
}

func NewTrigger() *Trigger {
	return &Trigger{}
}
