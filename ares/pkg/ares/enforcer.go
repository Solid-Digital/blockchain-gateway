package ares

import (
	"net/http"

	"bitbucket.org/unchain/ares/gen/orm"
)

type Enforcer interface {
	MakeUser(userID int64) error
	MakeSuperAdmin(userID int64) error
	IsSuperAdmin(userID int64) bool
	GetRolesForUserInOrganization(userID int64, organization string) map[string]bool
	GetGlobalRolesForUser(userID int64) map[string]bool
	SetMemberRoles(organizationName string, userID int64, roles map[string]bool) error
	GetAllRolesForUser(userID int64, orgs []*orm.Organization) map[string]map[string]bool
	Enforce(userID int64, path string, method string) bool
	Authorize(*http.Request, interface{}) error
	AddPolicy(params ...interface{}) bool
}
