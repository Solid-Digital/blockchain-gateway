package auth

import (
	"context"
	"database/sql"

	"github.com/go-openapi/strfmt"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"
	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) GetCurrentUser(u *dto.User) (*dto.GetCurrentUserResponse, *apperr.Error) {
	var user *orm.User

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		user, appErr = xorm.GetUserTxByID(ctx, tx, u.ID,
			qm.Load(orm.UserRels.Organizations),
			qm.Load(orm.UserRels.DefaultOrganization),
		)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	roles := s.enforcer.GetAllRolesForUser(user.ID, user.R.Organizations)

	defaultOrg, err := getDefaultOrg(user)
	if err != nil {
		return nil, apperr.Internal.Wrap(err)
	}

	ret := &dto.GetCurrentUserResponse{
		ID:                  user.ID,
		FullName:            user.FullName,
		Email:               user.Email.String,
		CreatedAt:           strfmt.DateTime(user.CreatedAt),
		DefaultOrganization: defaultOrg.Name,
		Roles:               roles,
	}

	return ret, nil
}
