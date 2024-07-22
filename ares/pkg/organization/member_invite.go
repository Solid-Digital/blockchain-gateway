package organization

import (
	"context"
	"database/sql"
	"errors"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/xorm"

	"github.com/volatiletech/null"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/ares"
)

func (s *Service) InviteMember(params *dto.InviteMemberRequest, orgName string) (*dto.InviteMemberResponse, *apperr.Error) {
	email := params.Email.String()
	roles := params.Roles

	var inviteID string
	var user *orm.User
	var member *dto.GetMemberResponse

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var err error
		var appErr *apperr.Error

		user, err = orm.Users(orm.UserWhere.Email.EQ(null.StringFrom(email))).One(ctx, tx)
		if err != nil {
			appErr = ares.ParsePQErr(err)
			switch {
			case errors.Is(appErr, apperr.NotFound):
				// If user doesn't exist, send them an invite
				inviteID, user, appErr = s.auth.InviteUserTx(ctx, tx, email, orgName)
				if appErr != nil {
					return appErr
				}
			default:
				return appErr
			}
		}

		member, appErr = s.InviteMemberTx(ctx, tx, user, orgName, roles)
		if appErr != nil {
			return appErr
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	ret := &dto.InviteMemberResponse{
		InviteID:  inviteID,
		ID:        member.ID,
		FullName:  member.FullName,
		Email:     member.Email,
		CreatedAt: member.CreatedAt,
		Roles:     member.Roles,
	}

	return ret, nil
}

func (s *Service) InviteMemberTx(ctx context.Context, tx *sql.Tx, user *orm.User, orgName string, roles map[string]bool) (*dto.GetMemberResponse, *apperr.Error) {
	var err error
	var appErr *apperr.Error

	org, appErr := xorm.GetOrganizationTx(ctx, tx, orgName)
	if appErr != nil {
		return nil, appErr
	}

	appErr = xorm.AddUserOrganizationTx(ctx, tx, user, org)
	if appErr != nil {
		return nil, appErr
	}

	if roles == nil {
		roles = make(map[string]bool)
	}
	roles[ares.RoleMember.String()] = true

	err = s.enforcer.SetMemberRoles(orgName, user.ID, roles)

	if err != nil {
		return nil, apperr.Internal.Wrap(err)
	}

	member := getMember(user, s.enforcer.GetRolesForUserInOrganization(user.ID, orgName))

	return member, nil
}
