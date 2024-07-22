package contract

import (
	"fmt"

	"github.com/unchainio/interfaces/logger"
	"github.com/unchainio/pkg/errors"
)

// Contracts maps contract addresses to Contract
type Contracts map[string]*Contract

func NewContracts(logger logger.Logger, cfgs []*Config) (Contracts, error) {
	contracts := Contracts{}

	for _, cfg := range cfgs {
		contract, err := newContract(logger, cfg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to configure contract")
		}

		contracts[cfg.Address] = contract
	}

	return contracts, nil
}

func (c Contracts) GetContract(address string) (*Contract, error) {
	// When a single Contract is configured, it is not required to specify that contract
	if address == "" && len(c) == 1 {
		for _, contract := range c {
			return contract, nil
		}
	}

	contract, ok := c[address]
	if !ok {
		return nil, errors.New(fmt.Sprintf("contract address %s invalid", address))
	}

	return contract, nil
}
