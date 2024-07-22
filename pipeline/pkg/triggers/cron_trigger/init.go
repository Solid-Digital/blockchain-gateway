package cron_trigger

import (
	"bytes"
	"github.com/robfig/cron/v3"
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchainio/pkg/errors"
	"github.com/unchainio/pkg/xconfig"
)

func (t *Trigger) Init(stub domain.Stub, config []byte) error {
	cfg := new(Config)
	err := xconfig.Load(cfg, xconfig.FromReaders("toml", bytes.NewReader(config)))
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}
	t.config = cfg
	t.stub = stub

	t.cron = cron.New(cron.WithSeconds())
	t.RequestChannel = make(chan *domain.Request)

	_, err = t.cron.AddFunc(cfg.Specification, func() {
		req := domain.NewRequest(nil)
		t.RequestChannel <-req
	})
	if err != nil {
		return err
	}

	// starts the cron in its own go routine
	t.cron.Start()

	return nil
}
