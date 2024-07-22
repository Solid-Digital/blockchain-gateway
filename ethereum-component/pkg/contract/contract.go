package contract

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/unchainio/interfaces/logger"
)

type Contract struct {
	logger  logger.Logger
	cfg     *Config
	Address *common.Address
	ABI     *abi.ABI
}

func newContract(logger logger.Logger, cfg *Config) (*Contract, error) {
	address := common.HexToAddress(cfg.Address)

	ABI, err := abi.JSON(strings.NewReader(cfg.ABI))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse ABI")
	}

	return &Contract{logger: logger, cfg: cfg, Address: &address, ABI: &ABI}, nil
}

func (c *Contract) GetFunction(functionName string) (*abi.Method, error) {
	function, ok := c.ABI.Methods[functionName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("function %s does not exist in ABI", functionName))
	}

	return &function, nil
}

func (c *Contract) GetEvent(eventName string) (*abi.Event, error) {
	event, ok := c.ABI.Events[eventName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("event %s does not exist in ABI", eventName))
	}

	return &event, nil
}

func (c *Contract) AddressString() string {
	return strings.ToLower(c.Address.String())
}
