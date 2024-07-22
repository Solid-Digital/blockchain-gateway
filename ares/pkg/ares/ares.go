package ares

import (
	"bitbucket.org/unchain/ares/gen/api"
	"bitbucket.org/unchain/ares/pkg/3p/sql"
	"github.com/unchainio/interfaces/logger"
)

// Server is the ares server struct
type Server struct {
	Handlers   *api.Handlers
	Middleware Middleware
	DB         *sql.DB
	Redis      KVStore
	FileStore  FileStore

	Enforcer Enforcer

	AuthService         AuthService
	ComponentService    ComponentService
	DeploymentService   DeploymentService
	OrganizationService OrganizationService
	PipelineService     PipelineService
	SubscriptionService SubscriptionService

	Log logger.Logger
}

// NewServer constructs a new Ares server
func NewServer(handlers *api.Handlers, middleware Middleware, db *sql.DB, redis KVStore, fileStore FileStore, enforcer Enforcer, authService AuthService, componentService ComponentService, deploymentService DeploymentService, organizationService OrganizationService, pipelineService PipelineService, subscriptionService SubscriptionService, log logger.Logger) (*Server, error) {
	server := &Server{
		Handlers:            handlers,
		Middleware:          middleware,
		DB:                  db,
		Redis:               redis,
		FileStore:           fileStore,
		Enforcer:            enforcer,
		AuthService:         authService,
		ComponentService:    componentService,
		DeploymentService:   deploymentService,
		OrganizationService: organizationService,
		PipelineService:     pipelineService,
		SubscriptionService: subscriptionService,
		Log:                 log,
	}

	return server, nil
}
