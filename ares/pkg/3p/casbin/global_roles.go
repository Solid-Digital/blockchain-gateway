package casbin

import "bitbucket.org/unchain/ares/pkg/ares"

func (e *Enforcer) MakeUser(userID int64) error {
	e.AddRoleForUser("*::"+idToString(userID), ares.RoleUser.String())

	return nil
}

func (e *Enforcer) MakeSuperAdmin(userID int64) error {
	e.AddRoleForUser("*::"+idToString(userID), ares.RoleSuperAdmin.String())

	return nil
}

func (e *Enforcer) IsSuperAdmin(userID int64) bool {
	roles := e.GetGlobalRolesForUser(userID)

	return roles[ares.RoleSuperAdmin.String()]
}

func (e *Enforcer) GetGlobalRolesForUser(userID int64) map[string]bool {
	return e.GetRolesForUserInOrganization(userID, "*")
}
