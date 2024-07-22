package api_trigger

import (
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchain/pipeline/pkg/triggers/api_trigger/config"
	"github.com/unchain/pipeline/pkg/triggers/api_trigger/listener"
)

const defaultPort = "80"

func (t *Trigger) Init(stub domain.Stub, cfg []byte) error {
	var err error

	t.cfg, err = config.NewConfig(cfg)
	if err != nil {
		return err
	}
	t.stub = stub

	port := defaultPort
	if t.cfg.Port != "" {
		port = t.cfg.Port
	}

	t.client, err = listener.NewServer(t.stub, t.cfg, port)
	if err != nil {
		return err
	}
	stub.Printf("http api trigger initialized")

	return nil
}
