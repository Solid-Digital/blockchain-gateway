package component

import (
	"bitbucket.org/unchain/ares/pkg/3p/sql"
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/google/wire"
	"github.com/unchainio/interfaces/logger"
)

var ServiceSet = wire.NewSet(NewService, wire.Bind(new(ares.ComponentService), new(Service)))

var _ ares.ComponentService = &Service{}

type Service struct {
	db *sql.DB

	service *service
	log     logger.Logger
}

type service struct {
	store    ares.FileStore
	enforcer ares.Enforcer
}

func NewService(
	db *sql.DB,
	store ares.FileStore,
	enforcer ares.Enforcer,
	log logger.Logger,
) (*Service, error) {
	return &Service{
		db: db,
		service: &service{
			store:    store,
			enforcer: enforcer,
		},
		log: log,
	}, nil
}
