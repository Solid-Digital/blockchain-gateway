package ethereum_client

import (
	"bitbucket.org/unchain/ethereum2/pkg/account"
	"bitbucket.org/unchain/ethereum2/pkg/contract"
	"github.com/bsm/redislock"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/unchainio/interfaces/logger"
)

type Client struct {
	logger    logger.Logger
	cfg       *Config
	redis     *redis.Client
	locker    *redislock.Client
	Client    *ethclient.Client
	Accounts  account.Accounts
	Contracts contract.Contracts
}

func NewClient(logger logger.Logger, cfg *Config, redis *redis.Client, accounts account.Accounts, contracts contract.Contracts) (*Client, error) {
	client, err := ethclient.Dial(cfg.Host)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start ethereum client")
	}

	locker := redislock.New(redis)

	return &Client{
		logger:    logger,
		cfg:       cfg,
		redis:     redis,
		locker:    locker,
		Client:    client,
		Accounts:  accounts,
		Contracts: contracts,
	}, nil
}
