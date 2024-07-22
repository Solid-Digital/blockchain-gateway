package action

import (
	"bitbucket.org/unchain/ethereum2/pkg/account"
	"bitbucket.org/unchain/ethereum2/pkg/contract"
	"bitbucket.org/unchain/ethereum2/pkg/ethereum_client"
	"bitbucket.org/unchain/ethereum2/pkg/redis"
	"github.com/unchainio/interfaces/adapter"
)

func (a *Action) Init(stub adapter.Stub, config []byte) error {
	var err error

	a.stub = stub

	a.cfg, err = newConfig(config)
	if err != nil {
		a.stub.Errorf(err.Error())

		return err
	}

	if a.cfg.Redis != nil {
		a.redis, err = redis.NewClient(a.cfg.Redis)
		if err != nil {
			a.stub.Errorf(err.Error())

			return err
		}
	}

	accounts, err := account.NewAccounts(stub, a.cfg.Accounts)
	if err != nil {
		a.stub.Errorf(err.Error())

		return err
	}

	contracts, err := contract.NewContracts(stub, a.cfg.Contracts)
	if err != nil {
		a.stub.Errorf(err.Error())

		return err
	}

	a.client, err = ethereum_client.NewClient(stub, a.cfg.Ethereum, a.redis, accounts, contracts)
	if err != nil {
		a.stub.Errorf(err.Error())

		return err
	}

	stub.Printf("ethereum action initialized for node %s", a.cfg.Ethereum.Host)
	if a.redis != nil {
		stub.Printf("redis support for ethereum action on %s", a.cfg.Redis.Host)
	}

	return nil
}
