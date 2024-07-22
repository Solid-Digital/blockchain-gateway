package organization

import (
	"bitbucket.org/unchain/ares/pkg/3p/harbor"
	"bitbucket.org/unchain/ares/pkg/3p/sql"
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(NewService, wire.Bind(new(ares.OrganizationService), new(Service)))

type Service struct {
	db       *sql.DB
	auth     ares.AuthService
	enforcer ares.Enforcer
	registry *harbor.Client
	mailer   ares.Mailer
}

func NewService(db *sql.DB, auth ares.AuthService, enforcer ares.Enforcer, mailer ares.Mailer, registry *harbor.Client) *Service {
	service := &Service{
		db:       db,
		auth:     auth,
		enforcer: enforcer,
		registry: registry,
		mailer:   mailer,
	}

	// dirty(?) hack to avoid initialization cycles (auth service has dependency on the organization service and vice versa)
	auth.SetOrganizationService(service)

	return service
}
