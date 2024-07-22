package auth

import (
	"time"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"
)

func (s *Service) GenerateToken(userID int64, email string, duration time.Duration) (string, *apperr.Error) {
	token, err := s.generateToken(userID, email, duration)
	if err != nil {
		return "", apperr.Internal.Wrap(err)
	}

	return token, nil
}
