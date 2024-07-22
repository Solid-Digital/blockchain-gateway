package organization

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/gen/orm"

	"bitbucket.org/unchain/ares/gen/dto"
)

func (s *Service) CreateOrganization(params *dto.CreateOrganizationRequest, principal *dto.User) (*dto.GetOrganizationResponse, *apperr.Error) {
	orgName := params.Name
	orgDisplayName := params.DisplayName

	var org *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error
		var user *orm.User

		user, appErr = xorm.GetUserTxByID(ctx, tx, principal.ID)
		if appErr != nil {
			return appErr
		}

		org, appErr = xorm.CreateOrganizationTx(ctx, tx, orgName, orgDisplayName, user)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	roles := ares.AllOrganizationRolesMap

	err := s.enforcer.SetMemberRoles(org.Name, principal.ID, roles)
	if err != nil {
		return nil, apperr.Internal.Wrap(err)
	}

	// If harbor was configured, upsert a project for this organization.
	if s.registry != nil {
		_, err = s.registry.UpsertProject(org.Name, false)
		if err != nil {
			return nil, apperr.Internal.Wrap(err)
		}
	}

	return getOrg(org), nil
}
