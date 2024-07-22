package domain

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type DeployContractResponse map[string]*DeployedContract

type DeployedContract struct {
	Address *common.Address
	Tx      *types.Transaction
	ABI     string
}

func (d DeployContractResponse) AddContract(address *common.Address, tx *types.Transaction, ABI string) {
	// For some reason the abi package stores addresses with uppercase characters. We translate everything to
	// lowercase since this is most common.
	addressKey := strings.ToLower(address.String())

	d[addressKey] = &DeployedContract{
		Address: address,
		Tx:      tx,
		ABI:     ABI,
	}
}
