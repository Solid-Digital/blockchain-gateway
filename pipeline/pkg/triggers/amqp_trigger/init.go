package amqp_trigger

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchainio/pkg/errors"
	"github.com/unchainio/pkg/xconfig"
	"sync"
)

func (t *Trigger) Init(stub domain.Stub, config []byte) error {
	cfg := new(Config)
	err := xconfig.Load(cfg, xconfig.FromReaders("toml", bytes.NewReader(config)))
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}
	t.config = cfg
	t.stub = stub

	t.RequestChannel = make(chan *domain.Request)

	t.ResponseChannelMap = new(sync.Map)

	if cfg.AmqpScheme == "" {
		cfg.AmqpScheme = "amqp://"
	}

	// start consuming messages
	conn, err := amqp.Dial(fmt.Sprintf("%s%s:%s@%s:%s", cfg.AmqpScheme, cfg.Username, cfg.Password, cfg.Domain, cfg.Port))
	if err != nil {
		return err
	}
	t.amqpConn = conn

	c, err := conn.Channel()
	if err != nil {
		return err
	}
	t.amqpChannel = c

	msgChannel, err := c.Consume(cfg.QueueName, cfg.ConsumerName, false, false, false, false, nil)
	if err != nil {
		return err
	}


	go func() {
		for msg := range msgChannel {
			req := &domain.Request{
				Tag: domain.NewTag(),
				Output: map[string]interface{}{
					"body": msg.Body,
					"delivery": msg,
				},
			}
			t.ResponseChannelMap.Store(req.Tag, msg)
			t.RequestChannel <- req
		}
	}()

	return nil
}
