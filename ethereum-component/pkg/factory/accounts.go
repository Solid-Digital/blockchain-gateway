package factory

import (
	"bitbucket.org/unchain/ethereum2/pkg/account"
)

func (f *Factory) AccountCfgs(cfgFile []byte) []*account.Config {
	return f.ActionCfg(cfgFile).Accounts
}

func (f *Factory) Accounts(cfgFile []byte) account.Accounts {
	cfg := f.AccountCfgs(cfgFile)
	accounts, err := account.NewAccounts(f.logger, cfg)

	f.suite.Require().NoError(err)

	return accounts
}
