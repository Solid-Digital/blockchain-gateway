package auth

import (
	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/pkg/3p/apperr"
)

func (s *Service) Logout(token *dto.Token) *apperr.Error {
	err := s.kv.BlacklistToken(token)
	if err != nil {
		return apperr.Internal.Wrap(err).WithMessage("logout failed")
	}

	return nil
}
