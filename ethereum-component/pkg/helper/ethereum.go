package helper

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"

	"bitbucket.org/unchain/ethereum2/pkg/contract"
	"bitbucket.org/unchain/ethereum2/pkg/domain"
	"github.com/ethereum/go-ethereum/accounts/abi"

	"bitbucket.org/unchain/ethereum2/pkg/ethereum_client"
	"github.com/ethereum/go-ethereum/common"
)

func (h *Helper) ContractAtAddress(ethereumClient *ethereum_client.Client, address *common.Address) bool {
	bytes, err := ethereumClient.Client.CodeAt(context.Background(), *address, nil)

	h.suite.Require().NoError(err)

	if len(bytes) == 0 {
		return false
	}

	return true
}

func (h *Helper) TransactionCommitted(ethereumClient *ethereum_client.Client, tx *types.Transaction) bool {
	_, pending, err := ethereumClient.Client.TransactionByHash(context.Background(), tx.Hash())

	h.suite.Require().NoError(err)

	return !pending
}

func (h *Helper) DeploySingleContract(ethereumClient *ethereum_client.Client, solidity string, constructorParams map[string]map[string]interface{}) (address string, stringABI string) {
	account := ethereumClient.Accounts[h.KeyFromMap(ethereumClient.Accounts)]

	msg := &domain.MessageDeployContract{
		From:              h.AddressToString(account.Address),
		Solidity:          solidity,
		ConstructorParams: constructorParams,
	}

	response, err := ethereumClient.DeployContracts(msg)

	h.suite.Require().NoError(err)

	deployedContract := response[h.KeyFromMap(response)]
	address = h.AddressToString(deployedContract.Address)
	ABI, err := abi.JSON(strings.NewReader(deployedContract.ABI))

	h.suite.Require().NoError(err)

	ethereumClient.Contracts[address] = &contract.Contract{
		Address: deployedContract.Address,
		ABI:     &ABI,
	}

	return address, stringABI
}
