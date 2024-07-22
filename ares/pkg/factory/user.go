package factory

import (
	"context"
	"database/sql"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/pkg/auth"

	"bitbucket.org/unchain/ares/pkg/testhelper"

	"time"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"github.com/Pallinder/go-randomdata"
	"github.com/go-openapi/strfmt"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

func (f *Factory) User(create bool) *orm.User {
	hashedPassword, err := auth.HashPassword("qwerty85")

	f.suite.NoError(err)

	user := &orm.User{
		FullName:              randomdata.FullName(randomdata.RandomGender),
		Email:                 null.StringFrom(testhelper.Randumb(randomdata.Email())),
		DefaultOrganizationID: null.NewInt64(0, false),
		PasswordHash:          hashedPassword,
	}

	if !create {
		return user
	}

	err = f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return user.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return user
}

func (f *Factory) AddUserToOrg(user *orm.User, org *orm.Organization) {
	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		err := org.AddUsers(ctx, tx, false, user)
		if err != nil {
			return err
		}

		return nil
	})

	f.suite.Require().NoError(err)
}

func (f *Factory) DTOUser(create bool) *dto.User {
	user := f.User(create)

	return f.ORMToDTOUser(user)
}

func (f *Factory) ORMToDTOUser(user *orm.User) *dto.User {
	return &dto.User{
		ID:    user.ID,
		Email: strfmt.Email(user.Email.String),
	}
}

func (f *Factory) InvitedUser(create bool) *orm.User {
	hashedPassword, err := auth.HashPassword("qwerty85")

	f.suite.NoError(err)

	user := &orm.User{
		Email:        null.StringFrom(testhelper.Randumb(randomdata.Email())),
		PasswordHash: hashedPassword,
	}

	if !create {
		return user
	}

	err = f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		return user.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return user
}

func (f *Factory) RegisteredUser() (*orm.User, *orm.AccountConfirmationToken) {
	var user *orm.User
	var ac *orm.AccountConfirmationToken
	var org *orm.Organization

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		org = f.Organization(true)

		user = f.InvitedUser(true)

		err := user.SetDefaultOrganization(ctx, tx, false, org)
		if err != nil {
			return err
		}
		user.Status = null.StringFrom(ares.StatusPendingConfirmation)
		_, err = user.Update(ctx, tx, boil.Infer())
		if err != nil {
			return err
		}

		ac = &orm.AccountConfirmationToken{
			Token:          randomdata.RandStringRunes(16),
			ExpirationTime: time.Now().UTC().Add(48 * time.Hour),
			UserID:         user.ID,
		}
		return ac.Insert(ctx, tx, boil.Infer())
	})

	f.suite.Require().NoError(err)

	return user, ac
}

// Time to think about refactoring the factories
func (f *Factory) RegisteredUser2() *orm.User {
	user, _ := f.RegisteredUser()

	return user
}

func (f *Factory) UserFromOrg(org *orm.Organization) *orm.User {
	var err error
	var user *orm.User

	err = f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		user, err = org.Users().One(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	f.suite.Require().NotNil(user)

	return user
}

func (f *Factory) ArchivedUser() *orm.User {
	user := f.User(true)
	appErr := f.ares.AuthService.DeleteCurrentUser(f.ORMToDTOUser(user))
	xrequire.NoError(f.suite.T(), appErr)

	return user
}
