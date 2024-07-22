package factory

import (
	"bitbucket.org/unchain/ethereum2/pkg/event"
)

func (f *Factory) EventCfgs(cfgFile []byte) []*event.Config {
	return f.TriggerCfg(cfgFile).Events
}
