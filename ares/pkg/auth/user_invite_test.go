package auth_test

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_InviteUserTx() {
	cases := map[string]struct {
		Email   string
		OrgName string
		Success bool
	}{
		"new email": {
			s.factory.User(false).Email.String,
			s.factory.Organization(true).Name,
			true},
		"email already exists": {
			s.factory.User(true).Email.String,
			s.factory.Organization(true).Name,
			false},
		"duplicate email with different case not allowed": {
			strings.ToUpper(s.factory.User(true).Email.String),
			s.factory.Organization(true).Name,
			false},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			var user *orm.User

			appErr := ares.WrapTx(s.ares.DB, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
				var appErr *apperr.Error

				_, user, appErr = s.service.InviteUserTx(ctx, tx, tc.Email, tc.OrgName)
				if appErr != nil {
					return appErr
				}

				return nil
			})

			if tc.Success {
				xrequire.NoError(t, appErr)
				require.Equal(t, tc.Email, user.Email.String)
				require.Equal(t, ares.StatusPendingConfirmation, user.Status.String)

				ac := s.helper.DBGetAccountConfirmation(user.ID)

				require.NotNil(t, ac.Token)
				require.True(t, ac.ExpirationTime.After(time.Now().Add(47*time.Hour)))
				require.True(t, ac.ExpirationTime.Before(time.Now().Add(49*time.Hour)))

				mail := s.helper.GetSignUpCodeFromEmail(ac.Token, user.Email.String)

				require.NotNil(t, mail)
			} else {
				xrequire.Error(t, appErr)
			}
		})
	}
}
