package http

import (
	"net/http"

	"github.com/go-chi/chi/middleware"

	"github.com/unchainio/pkg/xlogger"

	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/go-chi/chi"
)

func MiddlewareProvider(cfg *xlogger.Config) ares.Middleware {
	return chi.Chain(middleware.RequestID, middleware.RealIP, Logger, Recoverer(cfg.Level), CorsMiddleware).Handler
}

func InnerMiddlewareProvider() func(http.Handler) http.Handler {
	return nil
}
