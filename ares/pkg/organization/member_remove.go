package organization

import (
	"context"
	"database/sql"
	"errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) RemoveMember(email string, orgName string) *apperr.Error {
	return ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error
		var err error
		var user *orm.User
		var org *orm.Organization

		user, appErr = xorm.GetUserTxByEmail(ctx, tx, email)
		if appErr != nil {
			return appErr
		}

		org, appErr = xorm.GetUserOrganizationTx(ctx, tx, user, orgName)
		if appErr != nil {
			return appErr
		}

		err = user.RemoveOrganizations(ctx, tx, org)
		if err != nil {
			err := ares.ParsePQErr(err)
			switch {
			case errors.Is(err, apperr.NotFound):
				return ares.ErrUserEmailNotFound(err, email)
			default:
				return err
			}
		}

		err = s.enforcer.SetMemberRoles(orgName, user.ID, nil)
		if err != nil {
			return apperr.Internal.Wrap(err)
		}

		return nil
	})
}
