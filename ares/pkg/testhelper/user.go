package testhelper

import (
	"context"
	"database/sql"

	"github.com/volatiletech/null"

	"bitbucket.org/unchain/ares/gen/orm"
)

func (h *Helper) DBUserExists(userID int64) bool {
	var exists bool
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		exists, err = orm.Users(orm.UserWhere.ID.EQ(userID)).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return exists
}

func (h *Helper) MakeSuperAdmin(id int64) {
	err := h.ares.Enforcer.MakeSuperAdmin(id)
	h.suite.Require().NoError(err)
}

func (h *Helper) DBGetUser(userID int64) *orm.User {
	if !h.DBUserExists(userID) {
		return nil
	}

	var userFromDB *orm.User
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		userFromDB, err = orm.Users(orm.UserWhere.ID.EQ(userID)).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return userFromDB
}

func (h *Helper) DBGetUserByEmail(email string) *orm.User {
	var userFromDB *orm.User
	var err error

	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		userFromDB, err = orm.Users(orm.UserWhere.Email.EQ(null.StringFrom(email))).One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	h.suite.Require().NoError(err)

	return userFromDB
}

func (h *Helper) DBGetAccountConfirmation(userID int64) *orm.AccountConfirmationToken {
	var acFromDB *orm.AccountConfirmationToken
	var err error
	err = h.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		acFromDB, err = orm.AccountConfirmationTokens(orm.AccountConfirmationTokenWhere.UserID.EQ(userID)).One(ctx, tx)
		if err != nil {
			return err
		}
		return nil
	})

	h.suite.Require().NoError(err)

	return acFromDB
}
