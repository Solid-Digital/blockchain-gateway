package factory

import (
	"bitbucket.org/unchain/ethereum2/pkg/contract"
)

func (f *Factory) ContractCfgs(cfgFile []byte) []*contract.Config {
	return f.ActionCfg(cfgFile).Contracts
}

func (f *Factory) Contracts(cfgFile []byte) contract.Contracts {
	cfg := f.ContractCfgs(cfgFile)
	contracts, err := contract.NewContracts(f.logger, cfg)

	f.suite.Require().NoError(err)

	return contracts
}
