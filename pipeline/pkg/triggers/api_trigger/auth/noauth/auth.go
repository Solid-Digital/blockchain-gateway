package noauth

import "net/http"

type AuthService struct {
}

func NewService() *AuthService {
	return &AuthService{
	}
}

func (a *AuthService) Authenticate(r *http.Request) error {
	return nil
}
