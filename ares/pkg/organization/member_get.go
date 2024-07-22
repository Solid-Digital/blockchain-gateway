package organization

import (
	"context"
	"database/sql"
	"errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/ares"

	"github.com/volatiletech/null"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/go-openapi/strfmt"
)

func (s *Service) GetMember(email string, orgName string, principal *dto.User) (*dto.GetMemberResponse, *apperr.Error) {
	//var ret *dto.GetMemberResponse
	var org *orm.Organization
	var user *orm.User

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var err error
		var appErr *apperr.Error

		org, appErr = xorm.GetOrganizationTx(ctx, tx, orgName)
		if appErr != nil {
			return appErr
		}

		user, err = org.Users(orm.UserWhere.Email.EQ(null.StringFrom(email))).One(ctx, tx)
		if err != nil {
			err := ares.ParsePQErr(err)
			switch {
			case errors.Is(err, apperr.NotFound):
				return ares.ErrUserEmailNotFound(err, email)
			default:
				return err
			}
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	return getMember(user, s.enforcer.GetRolesForUserInOrganization(user.ID, orgName)), nil
}

func getMember(user *orm.User, roles map[string]bool) *dto.GetMemberResponse {
	return &dto.GetMemberResponse{
		CreatedAt: strfmt.DateTime(user.CreatedAt),
		Email:     user.Email.String,
		FullName:  user.FullName,
		ID:        user.ID,
		Roles:     roles,
	}
}
