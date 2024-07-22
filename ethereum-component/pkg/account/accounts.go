package account

import (
	"fmt"

	"github.com/unchainio/interfaces/logger"
	"github.com/unchainio/pkg/errors"
)

// Accounts maps account addresses to Account
type Accounts map[string]*Account

func NewAccounts(logger logger.Logger, cfgs []*Config) (Accounts, error) {
	accounts := Accounts{}

	for _, cfg := range cfgs {
		account, err := newAccount(logger, cfg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to configure account")
		}

		accounts[cfg.Address] = account
	}

	return accounts, nil
}

func (a Accounts) GetAccount(address string) (*Account, error) {
	// When a single Account is configured, it is not required to specify that account
	if address == "" && len(a) == 1 {
		for _, account := range a {
			return account, nil
		}
	}

	account, ok := a[address]
	if !ok {
		return nil, errors.New(fmt.Sprintf("account address %s invalid", address))
	}

	return account, nil
}
