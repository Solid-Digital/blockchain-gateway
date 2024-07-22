package subscription

import (
	"bitbucket.org/unchain/ares/pkg/3p/sql"
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/google/wire"
	"github.com/unchainio/interfaces/logger"
)

var ServiceSet = wire.NewSet(NewService, wire.Bind(new(ares.SubscriptionService), new(Service)))

// var ServiceSetOld = wire.NewSet()

var _ ares.SubscriptionService = &Service{}

type Service struct {
	db  *sql.DB
	AWS ares.AWSClient

	log logger.Logger
	cfg *Config
}

type Config struct {
	ActivateHandler bool
}

func NewService(
	cfg *Config,
	db *sql.DB,
	aws ares.AWSClient,
	log logger.Logger,
) (*Service, error) {
	s := &Service{
		cfg: cfg,
		AWS: aws,
		db:  db,
		log: log,
	}

	if cfg.ActivateHandler {
		s.handleAwsMarketplaceNotifications()
	}

	return s, nil
}
