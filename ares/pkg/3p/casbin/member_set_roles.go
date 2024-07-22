package casbin

import (
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/unchainio/pkg/errors"
)

func (e *Enforcer) SetMemberRoles(organizationName string, userID int64, roles map[string]bool) error {
	for role := range roles {
		if ok := ares.AllOrganizationRolesMap[role]; !ok {
			return errors.Errorf("role `%s` does not exist", role)
		}
	}

	e.DeleteRolesForUserInOrganization(userID, organizationName, ares.AllOrganizationRolesMap)

	// Add only the new roles
	e.AddRoleForUserInOrganization(userID, organizationName, roles)

	return nil
}
