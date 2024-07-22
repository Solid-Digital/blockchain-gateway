package casbin

import (
	"fmt"
	"time"

	"bitbucket.org/unchain/ares/pkg/3p/sql"
	"github.com/google/wire"

	"github.com/jinzhu/gorm"

	gormadapter "github.com/casbin/gorm-adapter"

	"bitbucket.org/unchain/ares/gen/orm"

	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/casbin/casbin"
	"github.com/unchainio/interfaces/logger"
)

var Set = wire.NewSet(NewEnforcer, wire.Bind(new(ares.Enforcer), new(Enforcer)))

type Enforcer struct {
	*casbin.SyncedEnforcer

	cfg *Config
	log logger.Logger
}

const Model = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(parseDomain(r.obj, p.obj) + "::" + r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
`

func NewEnforcer(log logger.Logger, db *sql.DB, cfg *Config) (*Enforcer, error) {
	m := casbin.NewModel(Model)

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		return nil, err
	}

	a := gormadapter.NewAdapterByDB(gormDB) // Your driver and data source.

	enforcer := casbin.NewSyncedEnforcer(m, a, false)
	enforcer.StartAutoLoadPolicy(30 * time.Second)
	enforcer.AddFunction("keyMatchDomain", KeyMatchDomainFunc)
	enforcer.AddFunction("parseDomain", ParseDomainFunc)

	return &Enforcer{
		SyncedEnforcer: enforcer,
		log:            log,
		cfg:            cfg,
	}, nil
}

func idToString(userID int64) string {
	return fmt.Sprintf("%d", userID)
}

func (e *Enforcer) Enforce(userID int64, path string, method string) bool {
	return e.Enforcer.Enforce(idToString(userID), path, method)
}

func (e *Enforcer) AddRoleForUserInOrganization(userID int64, organization string, roles map[string]bool) {
	for role, flag := range roles {
		if flag == true {
			e.AddRoleForUser(organization+"::"+idToString(userID), role)
		}
	}
}

func (e *Enforcer) GetAllRolesForUser(userID int64, orgs []*orm.Organization) map[string]map[string]bool {
	allRoles := make(map[string]map[string]bool)

	for _, org := range orgs {
		allRoles[org.Name] = e.GetRolesForUserInOrganization(userID, org.Name)
	}

	return allRoles
}

func (e *Enforcer) GetRolesForUserInOrganization(userID int64, organization string) map[string]bool {
	roleStrings := e.GetRolesForUser(organization + "::" + idToString(userID))

	rolesMap := make(map[string]bool)
	for _, roleString := range roleStrings {
		rolesMap[roleString] = true
	}

	return rolesMap
}

func (e *Enforcer) DeleteAllRolesForUserInOrganization(userID int64, organization string) {
	e.DeleteRolesForUser(organization + "::" + idToString(userID))
}

func (e *Enforcer) DeleteRolesForUserInOrganization(userID int64, organization string, roles map[string]bool) {
	for role, flag := range roles {
		if flag == true {
			e.DeleteRoleForUser(organization+"::"+idToString(userID), role)
		}
	}
}
