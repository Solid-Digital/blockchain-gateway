package xorm

import (
	"context"
	"database/sql"
	stderr "errors"
	"strings"
	"time"

	"github.com/volatiletech/sqlboiler/boil"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"github.com/volatiletech/null"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/ares"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

func GetUserTxByEmail(ctx context.Context, tx *sql.Tx, email string, mods ...qm.QueryMod) (*orm.User, *apperr.Error) {
	user, err := orm.Users(append(
		mods,
		orm.UserWhere.Email.EQ(null.StringFrom(email)))...,
	).One(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrUserEmailNotFound(err, email)
		default:
			return nil, err
		}
	}

	return user, nil
}

func GetUserTxByID(ctx context.Context, tx *sql.Tx, id int64, mods ...qm.QueryMod) (*orm.User, *apperr.Error) {
	user, err := orm.Users(append(
		mods,
		orm.UserWhere.ID.EQ(id))...,
	).One(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrUserIDNotFound(err, id)
		default:
			return nil, err
		}
	}

	return user, nil
}

func AddUserOrganizationTx(ctx context.Context, tx *sql.Tx, user *orm.User, org *orm.Organization) *apperr.Error {
	err := user.AddOrganizations(ctx, tx, false, org)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.Conflict):
			return ares.ErrDuplicateMembership(err, user.Email.String, org.Name)
		default:
			return err
		}
	}

	return nil
}

func GetUserOrganizationTx(ctx context.Context, tx *sql.Tx, user *orm.User, orgName string, mods ...qm.QueryMod) (*orm.Organization, *apperr.Error) {
	org, err := user.Organizations(append(
		mods,
		orm.OrganizationWhere.Name.EQ(orgName))...,
	).One(ctx, tx)
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.NotFound):
			return nil, ares.ErrNotMember(err, user.Email.String, orgName)
		default:
			return nil, err
		}
	}

	return org, nil
}

func ExistsUserOrganizationTx(ctx context.Context, tx *sql.Tx, user *orm.User, orgName string, mods ...qm.QueryMod) (bool, *apperr.Error) {
	exists, err := user.Organizations(append(
		mods,
		orm.OrganizationWhere.Name.EQ(orgName))...,
	).Exists(ctx, tx)

	if err != nil {
		return false, ares.ParsePQErr(err)
	}

	return exists, nil
}

func CreateUserTx(ctx context.Context, tx *sql.Tx, email string, token string) (*orm.User, *apperr.Error) {
	emailString := null.StringFrom(strings.ToLower(email))
	status := null.StringFrom(ares.StatusPendingConfirmation)

	user := &orm.User{
		Email:  emailString,
		Status: status,
	}

	err := user.Insert(ctx, tx, boil.Infer())
	if err != nil {
		err := ares.ParsePQErr(err)
		switch {
		case stderr.Is(err, apperr.Conflict):
			return nil, ares.ErrDuplicateUser(err, email)
		default:
			return nil, err
		}
	}

	ac := &orm.AccountConfirmationToken{
		Token:          token,
		UserID:         user.ID,
		ExpirationTime: time.Now().UTC().Add(48 * time.Hour),
	}

	err = ac.Insert(ctx, tx, boil.Infer())
	if err != nil {
		// Technically it would be possible to have a duplicate token error here, not
		// worth the effort to capture that.
		return nil, ares.ParsePQErr(err)
	}

	return user, nil
}
