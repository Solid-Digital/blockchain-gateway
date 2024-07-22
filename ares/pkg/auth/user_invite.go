package auth

import (
	"context"
	"database/sql"
	"fmt"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (s *Service) InviteUserTx(ctx context.Context, tx *sql.Tx, email string, orgName string) (inviteID string, user *orm.User, appErr *apperr.Error) {
	// create sign-up token
	inviteID, token, err := s.generateAuthCode()
	if err != nil {
		return "", nil, apperr.Internal.Wrap(err)
	}

	user, appErr = xorm.CreateUserTx(ctx, tx, email, token)
	if appErr != nil {
		return "", nil, appErr
	}

	err = s.enforcer.MakeUser(user.ID)
	if err != nil {
		return "", nil, apperr.Internal.Wrap(err)
	}

	email = user.Email.String // The case may be different from the initial value of email
	url := fmt.Sprintf(ares.URLFmt, s.cfg.ConnectURL, token, email, orgName, true)
	err = s.mailer.SendInviteMessage(email, inviteID, url, orgName)
	if err != nil {
		return "", nil, apperr.Internal.Wrap(err)
	}

	return inviteID, user, nil
}
