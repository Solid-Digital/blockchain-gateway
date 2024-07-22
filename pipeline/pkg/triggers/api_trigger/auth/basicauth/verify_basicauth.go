package basicauth

import (
	"errors"
	"net/http"
	"strings"
)

func (a *AuthService) Authenticate(r *http.Request) error {
	return a.verifyBasicAuth(r)
}

func (a *AuthService) verifyBasicAuth(r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if !strings.Contains(authHeader, "Basic") {
		return errors.New("auth header does not contain basicauth")
	}
	auth := strings.TrimPrefix(authHeader, "Basic ")

	if auth == "" {
		return errors.New("expected basicauth on request")
	}

	for _, a := range a.AuthStrings {
		if auth == a {
			return nil
		}
	}

	return errors.New("basic auth not valid")
}

