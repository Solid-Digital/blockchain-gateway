package listener

import (
	"fmt"
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchain/pipeline/pkg/triggers/api_trigger/auth"
	"github.com/unchain/pipeline/pkg/triggers/api_trigger/auth/apikey"
	"github.com/unchain/pipeline/pkg/triggers/api_trigger/auth/basicauth"
	"github.com/unchain/pipeline/pkg/triggers/api_trigger/auth/noauth"
	"github.com/unchain/pipeline/pkg/triggers/api_trigger/config"
	"net/http"
	"sync"

	"github.com/unchainio/interfaces/logger"
)

type Server struct {
	logger             logger.Logger
	cfg                *config.Config
	server             *http.Server
	auth               auth.Service
	RequestChannel     chan *domain.Request
	ResponseChannelMap *sync.Map
}

func NewServer(logger logger.Logger, cfg *config.Config, port string) (*Server, error) {
	authService := selectAuthService(cfg)

	server := &Server{
		logger:             logger,
		cfg:                cfg,
		server:             &http.Server{Addr: fmt.Sprintf(":%s", port)},
		auth:               authService,
		RequestChannel:     make(chan *domain.Request),
		ResponseChannelMap: new(sync.Map),
	}

	go server.runHTTPListener(server.server)

	return server, nil
}

func selectAuthService(cfg *config.Config) auth.Service {
	if cfg.Auth == nil {
		return noauth.NewService()
	}
	if len(cfg.Auth.APIKeys) > 0 {
		return apikey.NewService(cfg.Auth.APIKeys)
	}
	if len(cfg.Auth.BasicAuth) > 0 {
		return basicauth.NewService(cfg.Auth.BasicAuth)
	}
	return noauth.NewService()
}