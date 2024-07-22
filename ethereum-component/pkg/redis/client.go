package redis

import (
	"crypto/tls"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

func NewClient(cfg *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:      cfg.Host,
		Password:  cfg.Password,
		DB:        cfg.DB,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot connect to redis on: %s", cfg.Host)
	}

	return client, nil
}
