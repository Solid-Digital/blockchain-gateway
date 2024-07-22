package config

import (
	"github.com/BurntSushi/toml"
	"github.com/unchain/pipeline/pkg/triggers/api_trigger/auth/basicauth"
	"github.com/unchainio/pkg/errors"
)

type Auth struct {
	APIKeys        []string
	BasicAuth      []basicauth.Credentials
	AllowedOrigins []string
}

type Config struct {
	Port string `toml:"port"`
	Auth *Auth  `toml:"auth"`
}

func NewConfig(config []byte) (*Config, error) {
	c := new(Config)
	err := toml.Unmarshal(config, c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	if c.Auth != nil {
		err = c.validate()
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Config) validate() error {
	if len(c.Auth.APIKeys) > 0 && len(c.Auth.BasicAuth) > 0 {
		return errors.New("conflicting configuration, both API Keys and basic auth credentials configured")
	}

	if len(c.Auth.BasicAuth) > 0 {
		for _, b := range c.Auth.BasicAuth {
			if b.Username == "" || b.Password == "" {
				return errors.New("configuration error - could not find username or password for basic auth")
			}
		}
	}

	return nil
}
