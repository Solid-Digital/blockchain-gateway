package trigger

import (
	"bitbucket.org/unchain/ethereum2/pkg/contract"
	"bitbucket.org/unchain/ethereum2/pkg/ethereum_listener"
	"bitbucket.org/unchain/ethereum2/pkg/event"
	"github.com/unchainio/interfaces/adapter"
)

func (t *Trigger) Init(stub adapter.Stub, config []byte) error {
	var err error

	t.stub = stub

	t.cfg, err = newConfig(config)
	if err != nil {
		t.stub.Errorf(err.Error())

		return err
	}

	contracts, err := contract.NewContracts(t.stub, t.cfg.Contracts)
	if err != nil {
		t.stub.Errorf(err.Error())

		return err
	}

	events, err := event.NewEvents(t.stub, t.cfg.Events, contracts)
	if err != nil {
		t.stub.Errorf(err.Error())

		return err
	}

	listener, err := ethereum_listener.NewListener(t.stub, t.cfg.Ethereum, events)
	if err != nil {
		t.stub.Errorf(err.Error())

		return err
	}

	t.responseChannel = listener.ResponseChannel
	t.errorChannel = listener.ErrorChannel

	stub.Printf("ethereum trigger initialized for node %s", t.cfg.Ethereum.Host)

	return nil
}
