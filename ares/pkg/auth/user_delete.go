package auth

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"
	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"

	"bitbucket.org/unchain/ares/gen/dto"

	"bitbucket.org/unchain/ares/gen/orm"
)

const ArchivedFullName = "Archived user"

// When deleting a user, the email value will be set to nil and the names will be erased. This way,
// the user cannot login, but the FKs remain intact, without revealing personal details.
func (s *Service) DeleteCurrentUser(u *dto.User) *apperr.Error {
	return ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var user *orm.User
		var appErr *apperr.Error
		var err error

		user, appErr = xorm.GetUserTxByID(ctx, tx, u.ID)
		if appErr != nil {
			return appErr
		}

		user.Email = null.StringFromPtr(nil)
		user.FullName = ArchivedFullName

		_, err = user.Update(ctx, tx, boil.Whitelist(
			orm.UserColumns.Email,
			orm.UserColumns.FullName,
		))
		if err != nil {
			return ares.ParsePQErr(err).WithMessagef("failed to delete user %q", user.Email.String)
		}

		return nil

		// TODO(e-nikolov) make a postgres rule for soft deletes
		// TODO(e-nikolov) add business logic to check whether the user is allowed to be deleted (not the only member of an organization, etc..)
	})
}
