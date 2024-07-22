package auth

import (
	"context"
	"database/sql"
	"time"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"
	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (s *Service) ConfirmRegistration(params *dto.ConfirmRegistrationRequest) (*dto.LoginResponse, *apperr.Error) {
	password := params.Password
	fullName := params.FullName
	status := null.StringFrom(ares.StatusActive)

	var user *orm.User
	var hashedPw string

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		// Get registration information based on token
		ac, err := orm.AccountConfirmationTokens(orm.AccountConfirmationTokenWhere.Token.EQ(params.Token)).One(ctx, tx)
		if err != nil {
			return apperr.NotFound.WithMessage("token not found")
		}

		// Check if token is still valid
		if ac.ExpirationTime.Before(time.Now().UTC()) {
			return apperr.Forbidden.WithMessage("token is not valid anymore")
		}

		// Add fullname and password to user
		user, appErr = xorm.GetUserTxByID(ctx, tx, ac.UserID,
			qm.Load(orm.UserRels.Organizations),
			qm.Load(orm.UserRels.DefaultOrganization),
		)
		if appErr != nil {
			return appErr
		}

		hashedPw, appErr = s.HashPassword(password)
		if appErr != nil {
			return appErr
		}

		user.PasswordHash = hashedPw
		user.FullName = fullName
		user.Status = status

		_, err = user.Update(ctx, tx, boil.Infer())
		if err != nil {
			return ares.ParsePQErr(err)
		}

		_, err = ac.Delete(ctx, tx)
		if err != nil {
			return ares.ParsePQErr(err)
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	email := user.Email.String

	token, err := s.generateToken(user.ID, email, time.Hour*time.Duration(s.cfg.ExpirationDelta))
	if err != nil {
		return nil, apperr.Internal.Wrap(err)
	}

	defaultOrg, err := getDefaultOrg(user)
	if err != nil {
		return nil, apperr.Internal.Wrap(err)
	}

	ret := &dto.LoginResponse{
		DefaultOrganization: defaultOrg.Name,
		Token:               token,
		UserID:              user.ID,
	}

	return ret, nil
}
