package auth

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"github.com/go-chi/jwtauth"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"bitbucket.org/unchain/ares/gen/orm"

	"bitbucket.org/unchain/ares/gen/dto"

	"github.com/unchainio/pkg/errors"
)

func (s *Service) Login(ip string, params *dto.LoginRequest) (*dto.LoginResponse, *apperr.Error) {
	password := *params.Password
	email := strings.ToLower(params.Email.String())
	attempts, err := s.checkLoginAttempts(ip, email)
	if err != nil {
		return nil, apperr.Internal.Wrap(err)
	}

	ret, appErr := s.login(email, password, ip, attempts)
	// TODO if authentication error, increase login attempts here, not inside s.login() in order to keep only db access inside s.login() and no redis stuff
	if appErr != nil {
		return nil, appErr
	}

	s.kv.ClearLoginAttempts(email, ip, attempts)

	return ret, nil
}

func (s *Service) login(email string, password string, ip string, attempts int) (*dto.LoginResponse, *apperr.Error) {
	var user *orm.User
	var token string
	var defaultOrg *orm.Organization

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error
		var err error

		user, appErr = xorm.GetUserTxByEmail(ctx, tx, email,
			qm.Load(orm.UserRels.Organizations),
			qm.Load(orm.UserRels.DefaultOrganization),
		)
		if appErr != nil {
			return appErr
		}

		appErr = s.CompareHashAndPassword(user.PasswordHash, password)
		if appErr != nil {
			err := s.kv.IncrementLoginAttempts(ip, email, attempts)
			if err != nil {
				s.log.Errorf("failed to increase login attempt count - %v", err)
			}

			return appErr
		}

		token, err = s.generateToken(user.ID, user.Email.String, time.Hour*time.Duration(s.cfg.ExpirationDelta))
		if err != nil {
			return apperr.Internal.Wrap(err).WithMessage("Unable to generate token")
		}

		defaultOrg, err = getDefaultOrg(user)
		if err != nil {
			return apperr.Internal.Wrap(err)
		}

		return nil
	})
	if appErr != nil {
		return nil, appErr
	}

	ret := &dto.LoginResponse{
		DefaultOrganization: defaultOrg.Name,
		Token:               token,
		UserID:              user.ID,
	}

	return ret, nil
}

func getDefaultOrg(user *orm.User) (*orm.Organization, error) {
	if user.DefaultOrganizationID.Valid {
		return user.R.DefaultOrganization, nil
	}

	for _, org := range user.R.Organizations {
		return org, nil
	}

	return nil, errors.Errorf("user '%s' has no default organization", user.Email.String)
}

func (s *Service) checkLoginAttempts(ip, email string) (int, error) {
	attempts, err := s.kv.GetLoginAttempts(ip, email)
	if err != nil {
		s.log.Debugf("No unsuccessful login attempts since last login %s", email)
	}

	if attempts == 5 {
		// TODO specify the remaining time for the block rather than always saying 15 minutes
		return 0, errors.New("Your account was blocked for 15 minutes")
	}

	return attempts, nil
}

func (s *Service) generateToken(userID int64, email string, duration time.Duration) (string, error) {
	token, err := generateToken(s.cfg.Issuer, s.TokenAuth, userID, email, duration)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateToken(issuer string, tokenAuth *jwtauth.JWTAuth, userID int64, email string, duration time.Duration) (string, error) {
	timestamp := time.Now().UnixNano()
	expireToken := time.Now().Add(duration).Unix()

	idStr := strconv.FormatInt(userID, 10)
	claims := jwtauth.Claims{
		"id":    idStr,
		"email": email,
		"exp":   expireToken,
		"iss":   issuer,
		"date":  timestamp,
	}

	_, tokenString, err := tokenAuth.Encode(claims)

	if err != nil {
		return "", errors.Wrap(err, "Error getting signed token")
	}

	return tokenString, nil
}
