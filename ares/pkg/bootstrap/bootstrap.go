package bootstrap

import (
	"context"
	"fmt"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/volatiletech/null"

	"bitbucket.org/unchain/ares/pkg/ares"

	dbsql "database/sql"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/3p/sql"
	"github.com/unchainio/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
)

type Service struct {
	*service
	db  *sql.DB
	cfg *Config
}

type service struct {
	org      ares.OrganizationService
	auth     ares.AuthService
	enforcer ares.Enforcer
}

type Config struct {
	AdminEmailAddress string
	AdminPassword     string
	AdminOrganization string
}

// Bootstrap initially configures ares by creating an admin user
func New(db *sql.DB, organizationService ares.OrganizationService, enforcer ares.Enforcer, authService ares.AuthService, cfg *Config) *Service {
	return &Service{
		db: db,
		service: &service{
			org:      organizationService,
			auth:     authService,
			enforcer: enforcer,
		},
		cfg: cfg,
	}
}

func (s *Service) Bootstrap() *apperr.Error {
	fmt.Printf("Creating admin user...\n")

	return ares.WrapTx(s.db, func(ctx context.Context, tx *dbsql.Tx) *apperr.Error {
		var appErr *apperr.Error

		appErr = s.insertDefaultEnvironmentsTx(ctx, tx)
		if appErr != nil {
			return appErr
		}

		s.createCasbinRoles()

		admin, appErr := s.createAdminTx(ctx, tx)
		if appErr != nil {
			return appErr
		}

		adminOrg, appErr := s.createAdminOrgTx(ctx, tx, admin)
		if appErr != nil {
			return appErr
		}

		appErr = s.insertDefaultSubscriptionPlansTx(ctx, tx, admin)
		if appErr != nil {
			return appErr
		}

		appErr = s.insertDefaultBaseTx(ctx, tx, admin, adminOrg)
		if appErr != nil {
			return appErr
		}

		return nil
	})
}

func (s *Service) createAdminOrgTx(ctx context.Context, tx *dbsql.Tx, admin *orm.User) (*orm.Organization, *apperr.Error) {
	var appErr *apperr.Error

	orgExists, appErr := xorm.ExistsOrganizationTx(ctx, tx, s.cfg.AdminOrganization)
	if appErr != nil {
		return nil, appErr
	}

	var adminOrg *orm.Organization

	if !orgExists {
		adminOrg, appErr = xorm.CreateOrganizationTx(ctx, tx, s.cfg.AdminOrganization, s.cfg.AdminOrganization, admin)
		if appErr != nil {
			return nil, appErr.WithMessage("Bootstrap failed - Could not create Admin organization")
		}
	} else {
		adminOrg, appErr = xorm.GetOrganizationTx(ctx, tx, s.cfg.AdminOrganization)
		if appErr != nil {
			return nil, appErr
		}

		isMember, appErr := xorm.ExistsUserOrganizationTx(ctx, tx, admin, s.cfg.AdminOrganization)
		if appErr != nil {
			return nil, appErr
		}

		if !isMember {
			_, appErr = s.service.org.InviteMemberTx(ctx, tx, admin, s.cfg.AdminOrganization, nil)

			if appErr != nil {
				return nil, appErr.WithMessage("Bootstrap failed - Could not add admin to admin organization")
			}
		}
	}

	return adminOrg, nil
}

func (s *Service) createAdminTx(ctx context.Context, tx *dbsql.Tx) (*orm.User, *apperr.Error) {
	hashedPassword, appErr := s.service.auth.HashPassword(s.cfg.AdminPassword)
	if appErr != nil {
		return nil, appErr
	}

	admin := &orm.User{
		PasswordHash: hashedPassword,
		FullName:     "System Admin",
		Email:        null.StringFrom(s.cfg.AdminEmailAddress),
	}

	err := admin.Upsert(ctx, tx, true, []string{orm.UserColumns.Email}, boil.Infer(), boil.Infer())
	if err != nil {
		return nil, ares.ParsePQErr(errors.Wrap(err, "Bootstrap failed - Could not register Admin user"))
	}

	err = s.service.enforcer.MakeSuperAdmin(admin.ID)
	if err != nil {
		return nil, apperr.Internal.Wrap(err)
	}

	return admin, nil
}

func (s *Service) insertDefaultBaseTx(ctx context.Context, tx *dbsql.Tx, admin *orm.User, adminOrg *orm.Organization) *apperr.Error {
	base := &orm.Base{
		DeveloperID: adminOrg.ID,
		CreatedByID: admin.ID,
		UpdatedByID: admin.ID,
		Name:        "janus-v2",
		DisplayName: "Janus Base",
		Public:      true,
	}

	err := base.Upsert(ctx, tx, true, []string{orm.BaseColumns.Name}, boil.Infer(), boil.Infer())
	if err != nil {
		return ares.ParsePQErr(err)
	}

	baseVersion := &orm.BaseVersion{
		BaseID:         base.ID,
		CreatedByID:    admin.ID,
		UpdatedByID:    admin.ID,
		Version:        "v0.0.1",
		Description:    "Janus Base",
		Readme:         "",
		DockerImageRef: "registry.unchain.io/unchainio/janus-v2:latest",
		Entrypoint:     "janus",
		Public:         true,
	}

	err = baseVersion.Upsert(ctx, tx, true, []string{orm.BaseVersionColumns.BaseID, orm.BaseVersionColumns.Version}, boil.Infer(), boil.Infer())
	if err != nil {
		return ares.ParsePQErr(err)
	}

	return nil
}

func (s *Service) insertDefaultEnvironmentsTx(ctx context.Context, tx *dbsql.Tx) *apperr.Error {
	envs := []*orm.DefaultEnvironment{
		{
			Index: 1,
			Name:  "development",
		},
		{
			Index: 10000000,
			Name:  "production",
		},
	}

	for _, env := range envs {
		err := env.Upsert(ctx, tx, true, []string{orm.DefaultEnvironmentColumns.Name}, boil.Infer(), boil.Infer())
		if err != nil {
			return ares.ParsePQErr(err)
		}
	}

	return nil
}

func (s *Service) insertDefaultSubscriptionPlansTx(ctx context.Context, tx *dbsql.Tx, admin *orm.User) *apperr.Error {
	plans := []*orm.SubscriptionPlan{
		{
			CreatedByID:   admin.ID,
			UpdatedByID:   admin.ID,
			Name:          "Free Plan",
			PipelineLimit: 2,
		},
		{
			CreatedByID:   admin.ID,
			UpdatedByID:   admin.ID,
			Name:          "Starter Plan",
			PipelineLimit: 5,
		},
		{
			CreatedByID:   admin.ID,
			UpdatedByID:   admin.ID,
			Name:          "Professional Plan",
			PipelineLimit: 12,
		},
		{
			CreatedByID:   admin.ID,
			UpdatedByID:   admin.ID,
			Name:          "Enterprise Plan",
			PipelineLimit: 100,
		},
	}

	for _, plan := range plans {
		err := plan.Upsert(ctx, tx, true, []string{orm.DefaultEnvironmentColumns.Name}, boil.Infer(), boil.Infer())
		if err != nil {
			return ares.ParsePQErr(err)
		}
	}

	return nil
}

func (s *Service) createCasbinRoles() {
	// User - global role for users of the gateway. This role holds all permissions that aren't linked to a specific organization
	s.enforcer.AddPolicy(ares.RoleUser.String(), "/api/v1/auth", "(GET|POST|DELETE|PATCH|PUT)")
	s.enforcer.AddPolicy(ares.RoleUser.String(), "/api/v1/auth/*", "(GET|POST|DELETE|PATCH|PUT)")

	// SuperAdmin - global role for super admins of the gateway.
	s.enforcer.AddPolicy(ares.RoleSuperAdmin.String(), "/*", "(GET|POST|DELETE|PATCH|PUT)")

	// User Admin
	s.enforcer.AddPolicy(ares.RoleUserAdmin.String(), "/api/v1/orgs/:domain", "(GET|POST|DELETE|PATCH|PUT)")
	s.enforcer.AddPolicy(ares.RoleUserAdmin.String(), "/api/v1/orgs/:domain/members/*", "(GET|POST|DELETE|PATCH|PUT)")
	s.enforcer.AddPolicy(ares.RoleUserAdmin.String(), "/api/v1/orgs/:domain/members", "(GET|POST|DELETE|PATCH|PUT)")

	// Pipeline Operator
	s.enforcer.AddPolicy(ares.RolePipelineOperator.String(), "/api/v1/orgs/:domain/pipelines/*", "(POST|DELETE|PATCH|PUT)")
	s.enforcer.AddPolicy(ares.RolePipelineOperator.String(), "/api/v1/orgs/:domain/pipelines", "(POST|DELETE|PATCH|PUT)")
	// Component Developer
	s.enforcer.AddPolicy(ares.RoleComponentDeveloper.String(), "/api/v1/orgs/:domain/actions", "(POST|DELETE|PATCH|PUT)")
	s.enforcer.AddPolicy(ares.RoleComponentDeveloper.String(), "/api/v1/orgs/:domain/actions/*", "(POST|DELETE|PATCH|PUT)")
	s.enforcer.AddPolicy(ares.RoleComponentDeveloper.String(), "/api/v1/orgs/:domain/triggers", "(POST|DELETE|PATCH|PUT)")
	s.enforcer.AddPolicy(ares.RoleComponentDeveloper.String(), "/api/v1/orgs/:domain/triggers/*", "(POST|DELETE|PATCH|PUT)")

	// Member
	s.enforcer.AddPolicy(ares.RoleMember.String(), "/api/v1/orgs/:domain/pipelines/*", "(GET)")
	s.enforcer.AddPolicy(ares.RoleMember.String(), "/api/v1/orgs/:domain/pipelines", "(GET)")
	s.enforcer.AddPolicy(ares.RoleMember.String(), "/api/v1/orgs/:domain/actions/*", "(GET)")
	s.enforcer.AddPolicy(ares.RoleMember.String(), "/api/v1/orgs/:domain/actions", "(GET)")
	s.enforcer.AddPolicy(ares.RoleMember.String(), "/api/v1/orgs/:domain/triggers/*", "(GET)")
	s.enforcer.AddPolicy(ares.RoleMember.String(), "/api/v1/orgs/:domain/triggers", "(GET)")
	s.enforcer.AddPolicy(ares.RoleMember.String(), "/api/v1/orgs/:domain/bases/*", "(GET)")
	s.enforcer.AddPolicy(ares.RoleMember.String(), "/api/v1/orgs/:domain/bases", "(GET)")
}
