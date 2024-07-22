package factory

import (
	"bitbucket.org/unchain/ethereum2/pkg/redis"
	redis_client "github.com/go-redis/redis"
)

func (f *Factory) Redis(cfgFile []byte) *redis_client.Client {
	cfg := f.ActionCfg(cfgFile).Redis
	if cfg == nil {
		return nil
	}

	client, err := redis.NewClient(cfg)

	f.suite.Require().NoError(err)

	return client
}
