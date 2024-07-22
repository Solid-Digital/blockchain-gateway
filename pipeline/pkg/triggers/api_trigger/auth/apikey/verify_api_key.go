package apikey

import (
	"net/http"
	"strings"

	"github.com/unchainio/pkg/errors"
)


func (a *AuthService) Authenticate(r *http.Request) error {
	return a.verifyAPIKey(r)
}

func (a *AuthService) verifyAPIKey(req *http.Request) error {
	authHeader := req.Header.Get("Authorization")
	if !strings.Contains(authHeader, "Bearer") {
		return errors.New("auth header does not contain a Bearer token")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// If the token is empty...
	if token == "" {
		return errors.New("no api key provided")
	}

	// Check token is valid
	for _, key := range a.APIKeys {
		if token == key {
			return nil
		}
	}

	// Token is invalid
	return errors.New("api key not valid")
}
