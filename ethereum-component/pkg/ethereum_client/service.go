package ethereum_client

import (
	"math/big"

	"bitbucket.org/unchain/ethereum2/pkg/domain"
	"github.com/unchainio/pkg/errors"
)

func (c *Client) CallContractFunction(msg *domain.MessageCallContractFunction) (interface{}, error) {
	account, err := c.Accounts.GetAccount(msg.From)
	if err != nil {
		return nil, err
	}

	contract, err := c.Contracts.GetContract(msg.To)
	if err != nil {
		return nil, err
	}

	function, err := contract.GetFunction(msg.Function)
	if err != nil {
		return nil, err
	}

	params, err := functionInput(function, msg.Params)
	if err != nil {
		return nil, err
	}

	var response interface{}
	switch c.functionType(function) {
	case Call:
		response, err = c.call(account, contract, function, params)
	case Transaction:
		response, err = c.transaction(account, contract, function, params, msg.Nonce)
	default:
		return nil, errors.New("invalid function type")
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) DeployContracts(msg *domain.MessageDeployContract) (domain.DeployContractResponse, error) {
	account, err := c.Accounts.GetAccount(msg.From)
	if err != nil {
		return nil, err
	}

	solidity := msg.Solidity

	params := msg.ConstructorParams

	var nonce *big.Int
	//Use nonce provided by the caller
	if msg.Nonce > 0 {
		nonce = big.NewInt(int64(msg.Nonce))
	} else if c.redis != nil {
		// Without redis we don't do nonce management
		nonce, err = c.getNonce(account)
		if err != nil {
			return nil, err
		}
	}

	response, err := c.deploy(account, solidity, params, nonce)
	if err != nil {
		return nil, err
	}

	return response, nil
}
