package ethereum_client

import (
	"bitbucket.org/unchain/ethereum2/pkg/account"
	"bitbucket.org/unchain/ethereum2/pkg/contract"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/unchainio/pkg/errors"
)

func (c *Client) call(account *account.Account, contract *contract.Contract, function *abi.Method, params []interface{}) (interface{}, error) {
	boundContract := c.boundContract(contract)
	opts := c.callOptions(account)

	//TODO: try if you can take the response var out of this switch/case statement
	//TODO: try make interface -> []interface{} with length 1
	switch callOutputType(function) {
	case TypeSingleOutput:
		ret := singleOutput(function)
		err := c.singleOutput(boundContract, opts, &ret, function.Name, params)
		if err != nil {
			return nil, err
		}

		return ret, nil
	case TypeMultipleOutput:
		ret := multipleOutput(function)
		err := c.multipleOutput(boundContract, opts, ret, function.Name, params)
		if err != nil {
			return nil, err
		}

		return ret, nil
	default:
		return nil, errors.New("invalid output type")
	}
}

func (c *Client) callOptions(account *account.Account) *bind.CallOpts {
	// TODO: implement gas, gas limit and maybe context?
	return &bind.CallOpts{
		From: *account.Address,
	}
}

func (c *Client) singleOutput(contract *bind.BoundContract, opts *bind.CallOpts, out interface{}, function string, params []interface{}) error {
	err := contract.Call(opts, out, function, params...)
	if err != nil {
		return errors.Wrapf(err, "failed to call function %s", function)
	}

	return nil
}

func (c *Client) multipleOutput(contract *bind.BoundContract, opts *bind.CallOpts, out []interface{}, function string, params []interface{}) error {
	err := contract.Call(opts, &out, function, params...)
	if err != nil {
		return errors.Wrapf(err, "failed to call function %s", function)
	}

	return nil
}
