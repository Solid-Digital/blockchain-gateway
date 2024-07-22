package ares

import (
	"github.com/unchainio/pkg/errors"
)

type Role int

func ParseRole(roleString string) (string, error) {
	for _, role := range AllOrganizationRoles {
		if role == roleString {
			return role, nil
		}
	}

	return "", errors.Errorf("unknown role")
}

func (r Role) String() string {
	roles := [...]string{
		"SuperAdmin",
		"UserAdmin",
		"ComponentDeveloper",
		"PipelineOperator",
		"Member",
		"User",
	}

	if r < 0 || int(r) >= len(roles) {
		return "unknown"
	}

	return roles[r]
}

const (
	RoleSuperAdmin Role = iota
	RoleUserAdmin
	RoleComponentDeveloper
	RolePipelineOperator
	RoleMember
	RoleUser
)

func init() {
	AllOrganizationRolesMap = make(map[string]bool)

	for _, elem := range AllOrganizationRoles {
		AllOrganizationRolesMap[elem] = true
	}
}

var AllOrganizationRoles = []string{
	RoleUserAdmin.String(),
	RoleComponentDeveloper.String(),
	RolePipelineOperator.String(),
	RoleMember.String(),
}

var AllOrganizationRolesMap map[string]bool
