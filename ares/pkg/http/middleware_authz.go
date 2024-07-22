package http

import (
	"net/http"

	"bitbucket.org/unchain/ares/pkg/ares"
)

func AuthorizationMiddleware(enforcer ares.Enforcer) func(*http.Request, interface{}) error {
	return enforcer.Authorize
}
