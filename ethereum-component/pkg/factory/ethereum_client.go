package factory

import (
	"bitbucket.org/unchain/ethereum2/pkg/ethereum_client"
)

func (f *Factory) EthereumClientCfg(cfgFile []byte) *ethereum_client.Config {
	return f.ActionCfg(cfgFile).Ethereum
}

func (f *Factory) EthereumClient(cfgFile []byte) *ethereum_client.Client {
	ethereumClientCfg := f.EthereumClientCfg(cfgFile)

	accounts := f.Accounts(cfgFile)

	contracts := f.Contracts(cfgFile)

	redis := f.Redis(cfgFile)

	client, err := ethereum_client.NewClient(f.logger, ethereumClientCfg, redis, accounts, contracts)

	f.suite.Require().NoError(err)

	return client
}
