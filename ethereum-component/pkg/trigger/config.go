package trigger

import (
	"bitbucket.org/unchain/ethereum2/pkg/contract"
	"bitbucket.org/unchain/ethereum2/pkg/ethereum_listener"
	"bitbucket.org/unchain/ethereum2/pkg/event"
	"github.com/BurntSushi/toml"
	"github.com/unchainio/pkg/errors"
)

type Config struct {
	Ethereum  *ethereum_listener.Config
	Contracts []*contract.Config
	Events    []*event.Config
}

func newConfig(config []byte) (*Config, error) {
	cfg := new(Config)
	err := toml.Unmarshal(config, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	err = cfg.validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Consider using validator on struct instead:
// https://github.com/go-playground/validator
func (c *Config) validate() error {
	if c.Ethereum == nil {
		return errors.New("ethereum host not configured")
	}

	return nil
}
