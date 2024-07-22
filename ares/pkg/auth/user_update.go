package auth

import (
	"context"
	"database/sql"
	"strings"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/go-openapi/strfmt"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (s *Service) UpdateCurrentUser(params *dto.UpdateCurrentUserRequest, u *dto.User) (*dto.GetCurrentUserResponse, *apperr.Error) {
	email := null.StringFrom(strings.ToLower(string(params.Email)))
	fullName := params.FullName
	updatedByID := null.Int64From(u.ID)

	var user *orm.User

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		user, appErr = xorm.GetUserTxByID(ctx, tx, u.ID,
			qm.Load(orm.UserRels.Organizations))
		if appErr != nil {
			return appErr
		}

		if params.Email != "" {
			user.Email = email
		}
		if params.FullName != "" {
			user.FullName = fullName
		}

		user.UpdatedByID = updatedByID

		_, err := user.Update(ctx, tx, boil.Whitelist(
			orm.UserColumns.Email,
			orm.UserColumns.FullName,
			orm.UserColumns.UpdatedByID,
		))
		if err != nil {
			return ares.ParsePQErr(err).WithMessagef("failed to update user %q", email.String)
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	ret := &dto.GetCurrentUserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email.String,
		CreatedAt: strfmt.DateTime(user.CreatedAt),
		Roles:     s.enforcer.GetAllRolesForUser(user.ID, user.R.Organizations),
	}

	return ret, nil
}
