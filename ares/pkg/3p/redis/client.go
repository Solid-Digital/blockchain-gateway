package redis

import (
	"time"

	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/go-redis/redis"
	"github.com/google/wire"
	"github.com/unchainio/interfaces/logger"
	"github.com/unchainio/pkg/errors"
	"github.com/unchainio/pkg/iferr"
)

var ClientSet = wire.NewSet(NewClient, wire.Bind(new(ares.KVStore), new(Client)))

var _ ares.KVStore = &Client{}

// Client
type Client struct {
	client *redis.Client
	log    logger.Logger
	cfg    *Config
}

func NewClient(log logger.Logger, cfg *Config) (*Client, func(), error) {
	rc := redis.NewClient(&redis.Options{
		Addr:        cfg.Host,
		Password:    cfg.AuthToken,
		DialTimeout: 100 * time.Second,
		IdleTimeout: 240 * time.Second,
	})

	client := &Client{
		log:    log,
		cfg:    cfg,
		client: rc,
	}

	cleanup := func() {
		iferr.Warn(rc.Close())
	}

	return client, cleanup, nil
}

func (c *Client) Set(key string, value string, expiration time.Duration) error {
	_, err := c.client.Set(key, value, expiration).Result()

	if err != nil {
		return errors.Wrap(err, "failed to set value in redis")
	}

	return nil
}

func (c *Client) Get(key string) (string, error) {
	val, err := c.client.Get(key).Result()

	//c.log.Debugf("redis key: %s\n", key)

	if err != nil {
		return "", errors.Wrap(err, "failed to get value from redis")
	}

	return val, nil
}

func (c *Client) Delete(keys ...string) (int64, error) {
	val, err := c.client.Del(keys...).Result()
	if err != nil {
		return 0, errors.Wrap(err, "failed to remove keys")
	}

	return val, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}
