package http

import (
	"bitbucket.org/unchain/ares/pkg/ares"

	dto "bitbucket.org/unchain/ares/gen/dto"
)

func AuthenticationMiddleware(s ares.AuthService) func(token string) (*dto.User, error) {
	return s.Authenticate
}
