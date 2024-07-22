package organization

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/gen/dto"
)

func (s *Service) SetMemberRoles(params *dto.SetMemberRolesRequest, email string, orgName string, u *dto.User) *apperr.Error {
	return ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		user, appErr := xorm.GetUserTxByEmail(ctx, tx, email)
		if appErr != nil {
			return appErr
		}

		if u.ID == user.ID {
			return apperr.BadRequest.WithMessage("Cannot change your own roles, please ask another admin")
		}

		var exists bool
		exists, appErr = xorm.ExistsUserOrganizationTx(ctx, tx, user, orgName)
		if appErr != nil {
			return appErr
		}

		if !exists {
			return apperr.BadRequest.WithMessagef("User %s is not a member of organization %s", email, orgName)
		}

		params.Roles[ares.RoleMember.String()] = true

		err := s.enforcer.SetMemberRoles(orgName, user.ID, params.Roles)
		if err != nil {
			return apperr.Internal.Wrap(err)
		}

		return nil
	})
}
