package auth

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"
	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/auth/internal"
	"github.com/unchainio/pkg/xrand"

	"bitbucket.org/unchain/ares/gen/orm"

	"bitbucket.org/unchain/ares/gen/dto"
)

func (s *Service) ResetPassword(params *dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, *apperr.Error) {
	email := string(params.Email)

	var user *orm.User
	var reqID string

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error
		var err error
		var code string

		user, appErr = xorm.GetUserTxByEmail(ctx, tx, email)
		if appErr != nil {
			return appErr
		}

		reqID, code, err = s.generateAuthCode()
		if err != nil {
			return apperr.Internal.Wrap(err)
		}

		err = s.kv.StoreRecoveryCode(code, email)
		if err != nil {
			return apperr.Internal.Wrap(err)
		}

		//send email with code to user
		err = s.mailer.SendRecoveryMessage(email, reqID, s.cfg.ConnectURL, code, user.FullName)
		if err != nil {
			return apperr.Internal.Wrap(err)
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	ret := &dto.ResetPasswordResponse{
		RequestID: reqID,
	}

	return ret, nil
}

func (s *Service) generateAuthCode() (id string, code string, err error) {
	pwd, err := internal.GeneratePassword(100, 5, 0, false, true)
	if err != nil {
		return "", "", err
	}

	return xrand.String(15), pwd, nil
}
