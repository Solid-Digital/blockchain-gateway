package event

import (
	"bitbucket.org/unchain/ethereum2/pkg/contract"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/unchainio/interfaces/logger"
)

type Event struct {
	Contract *contract.Contract
	Event    *abi.Event
	Filters  [][]interface{}
}

func newEvent(logger logger.Logger, cfg *Config, contracts contract.Contracts) (*Event, error) {
	contract, err := contracts.GetContract(cfg.ContractAddress)
	if err != nil {
		return nil, err
	}

	event, err := contract.GetEvent(cfg.Name)
	if err != nil {
		return nil, err
	}

	filters, err := buildFilters(event.Inputs, cfg.Filters)
	if err != nil {
		return nil, err
	}

	return &Event{
		Contract: contract,
		Event:    event,
		Filters:  filters,
	}, nil
}
