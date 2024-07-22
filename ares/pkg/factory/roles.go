package factory

import "bitbucket.org/unchain/ares/pkg/ares"

func (f *Factory) DefaultRoles() map[string]bool {
	return map[string]bool{ares.RoleMember.String(): true}
}

func (f *Factory) SomeRoles() map[string]bool {
	return map[string]bool{
		ares.RoleMember.String():           true,
		ares.RolePipelineOperator.String(): true,
		ares.RoleUserAdmin.String():        true,
	}
}

func (f *Factory) AllRoles() map[string]bool {
	var roles = make(map[string]bool)
	for _, role := range ares.AllOrganizationRoles {
		roles[role] = true
	}

	return roles
}
