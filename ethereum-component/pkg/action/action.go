package action

import (
	"bitbucket.org/unchain/ethereum2/pkg/ethereum_client"
	"github.com/go-redis/redis"
	"github.com/unchainio/interfaces/adapter"
)

type Action struct {
	cfg    *Config
	stub   adapter.Stub
	client *ethereum_client.Client
	redis  *redis.Client
}
