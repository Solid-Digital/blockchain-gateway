package action

import (
	"bitbucket.org/unchain/ethereum2/pkg/account"
	"bitbucket.org/unchain/ethereum2/pkg/contract"
	"bitbucket.org/unchain/ethereum2/pkg/ethereum_client"
	"bitbucket.org/unchain/ethereum2/pkg/redis"
	"github.com/pelletier/go-toml"
	"github.com/unchainio/pkg/errors"
)

type Config struct {
	Accounts  []*account.Config
	Contracts []*contract.Config
	Ethereum  *ethereum_client.Config
	Redis     *redis.Config
}

func newConfig(config []byte) (*Config, error) {
	cfg := new(Config)
	err := toml.Unmarshal(config, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "configuration file not valid")
	}

	err = cfg.validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if len(c.Accounts) == 0 {
		return errors.New("no accounts configured")
	}

	if c.Ethereum == nil {
		return errors.New("ethereum host not configured")
	}

	return nil
}
