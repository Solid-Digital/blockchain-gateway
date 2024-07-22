package ethereum_client

import (
	"math/big"
	"time"

	"github.com/bsm/redislock"

	"bitbucket.org/unchain/ethereum2/pkg/account"
	"bitbucket.org/unchain/ethereum2/pkg/contract"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Client) transaction(account *account.Account, contract *contract.Contract, function *abi.Method, params []interface{}, userSuppliedNonce uint64) (*types.Transaction, error) {
	var nonce *big.Int
	var err error

	//Use nonce provided by the caller
	if userSuppliedNonce > 0 {
		nonce = big.NewInt(int64(userSuppliedNonce))
	} else if c.redis != nil {
		lock, err := c.locker.Obtain(getLockKey(account), 2*time.Minute, &redislock.Options{
			RetryStrategy: redislock.ExponentialBackoff(1*time.Second, 8*time.Second),
			Metadata:      "",
			Context:       nil,
		})
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		defer func() {
			err := lock.Release()
			if err != nil {
				c.logger.Errorf("%+v", err)
			}
		}()

		// Without redis we don't do nonce management
		nonce, err = c.getNonce(account)
		if err != nil {
			return nil, err
		}
	}

	boundContract := c.boundContract(contract)
	opts, err := c.transactionOptions(account, nonce, c.cfg.GasPrice, c.cfg.GasLimit)
	if err != nil {
		return nil, err
	}

	tx, err := boundContract.Transact(opts, function.Name, params...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to call function %s", function.Name)
	}

	//Wait for tx to be committed
	err = c.txCommitted(tx.Hash())
	if err != nil {
		return nil, err
	}

	//Without redis we don't do nonce management
	if c.redis != nil {
		//Keep track of latest nonce
		err = c.setNonce(account, tx.Nonce())
		if err != nil {
			//We don't error since the transaction has been sent to the node successfully
			c.logger.Debugf(err.Error())
		}
	}

	return tx, nil
}

func (c *Client) transactionOptions(account *account.Account, nonce *big.Int, gasPrice *int64, gasLimit uint64) (*bind.TransactOpts, error) {
	// TODO: implement gas, gas limit and maybe context?
	signerFn, err := account.GetSignerFn()
	if err != nil {
		return nil, err
	}

	var gasPriceOpt *big.Int

	if gasPrice != nil {
		gasPriceOpt = big.NewInt(*gasPrice)
	}

	return &bind.TransactOpts{
		From:     *account.Address,
		Signer:   signerFn,
		Nonce:    nonce,
		GasPrice: gasPriceOpt,
		GasLimit: gasLimit,
	}, nil
}
