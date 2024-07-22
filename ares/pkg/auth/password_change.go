package auth

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"
	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/volatiletech/sqlboiler/boil"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) ChangeCurrentPassword(params *dto.ChangeCurrentPasswordRequest, u *dto.User) *apperr.Error {
	return ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		user, appErr := xorm.GetUserTxByID(ctx, tx, u.ID)
		if appErr != nil {
			return appErr
		}

		appErr = s.CompareHashAndPassword(user.PasswordHash, params.CurrentPassword)
		if appErr != nil {
			return appErr
		}

		user.PasswordHash, appErr = s.HashPassword(string(params.NewPassword))
		if appErr != nil {
			return appErr
		}

		_, err := user.Update(ctx, tx, boil.Whitelist(orm.UserColumns.PasswordHash))
		if err != nil {
			return ares.ParsePQErr(err)
		}

		return nil
	})
}
