package auth

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"
	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/volatiletech/sqlboiler/boil"

	"bitbucket.org/unchain/ares/gen/orm"

	"bitbucket.org/unchain/ares/gen/dto"
)

func (s *Service) ConfirmResetPassword(params *dto.ConfirmResetPasswordRequest) *apperr.Error {
	return ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		email, err := s.kv.GetEmailByRecoveryCode(params.RecoveryCode)
		if err != nil {
			return apperr.Internal.Wrap(err)
		}

		user, appErr := xorm.GetUserTxByEmail(ctx, tx, email)
		if appErr != nil {
			return appErr
		}

		user.PasswordHash, appErr = s.HashPassword(params.Password)
		if appErr != nil {
			return appErr
		}

		_, err = user.Update(ctx, tx, boil.Whitelist(orm.UserColumns.PasswordHash))
		if err != nil {
			return ares.ParsePQErr(err)
		}

		err = s.kv.RemoveRecoveryCode(params.RecoveryCode)
		if err != nil {
			return apperr.Internal.Wrap(err)
		}

		return nil
	})
}
