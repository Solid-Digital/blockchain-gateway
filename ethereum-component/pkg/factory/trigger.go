package factory

import (
	"bitbucket.org/unchain/ethereum2/pkg/trigger"
	"github.com/BurntSushi/toml"
)

func (f *Factory) TriggerCfg(cfgFile []byte) *trigger.Config {
	cfg := new(trigger.Config)
	err := toml.Unmarshal(cfgFile, cfg)

	f.suite.Require().NoError(err)

	return cfg
}

func (f *Factory) InitializedTrigger(cfg []byte) *trigger.Trigger {
	trigger := new(trigger.Trigger)
	trigger.Init(f.logger, cfg)

	return trigger
}
