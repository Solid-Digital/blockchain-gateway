package factory

import (
	"bitbucket.org/unchain/ethereum2/pkg/action"
	"github.com/BurntSushi/toml"
)

func (f *Factory) ActionCfg(cfgFile []byte) *action.Config {
	cfg := new(action.Config)
	err := toml.Unmarshal(cfgFile, cfg)

	f.suite.Require().NoError(err)

	return cfg
}

func (f *Factory) InitializedAction(cfgFile []byte) *action.Action {
	action := new(action.Action)
	action.Init(f.logger, cfgFile)

	return action
}
