package auth

import (
	"bitbucket.org/unchain/ares/pkg/3p/mail"
	"bitbucket.org/unchain/ares/pkg/3p/sql"
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/go-chi/jwtauth"
	"github.com/google/wire"
	"github.com/unchainio/interfaces/logger"
)

var ServiceSet = wire.NewSet(NewService, wire.Bind(new(ares.AuthService), new(Service)))

// var ServiceSetOld = wire.NewSet()

var _ ares.AuthService = &Service{}

type Service struct {
	kv                  ares.KVStore
	db                  *sql.DB
	TokenAuth           *jwtauth.JWTAuth
	enforcer            ares.Enforcer
	organizationService ares.OrganizationService
	mailer              *mail.Mailer
	AWS                 ares.AWSClient

	log logger.Logger
	cfg *Config
}

func NewService(
	db *sql.DB,
	enforcer ares.Enforcer,
	mailer *mail.Mailer,
	client ares.KVStore,
	aws ares.AWSClient,
	log logger.Logger,
	cfg *Config,
) (*Service, error) {
	privateKey, err := cfg.TLS.Key.ParsePKCS1PrivateKey()

	if err != nil {
		return nil, err
	}

	publicKey, err := cfg.TLS.Cert.ParsePKCS1PublicKey2()

	if err != nil {
		return nil, err
	}

	tokenAuth := jwtauth.New("RS512", privateKey, publicKey)

	s := &Service{
		kv:        client,
		db:        db,
		TokenAuth: tokenAuth,
		enforcer:  enforcer,
		mailer:    mailer,
		AWS:       aws,
		log:       log,
		cfg:       cfg,
	}

	return s, nil
}

func (s *Service) SetOrganizationService(service ares.OrganizationService) {
	s.organizationService = service
}
