package auth

import (
	"bitbucket.org/unchain/ares/pkg/3p/apperr"
	"github.com/unchainio/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) HashPassword(password string) (string, *apperr.Error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", apperr.Internal.Wrap(err)
	}

	return hashedPassword, nil
}

func (s *Service) CompareHashAndPassword(hash, password string) *apperr.Error {
	err := CompareHashAndPassword(hash, password)
	if err != nil {
		return apperr.Internal.Wrap(err)
	}

	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate password")
	}

	return string(hashedPassword), nil
}

func CompareHashAndPassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.Wrap(err, "incorrect password")
	}

	return nil
}
