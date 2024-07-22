package listener

import (
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchain/pipeline/pkg/triggers/api_trigger/output"
	"github.com/unchainio/pkg/errors"
	"github.com/unchainio/pkg/iferr"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func (s *Server) runHTTPListener(server *http.Server) {
	r := chi.NewRouter()
	corsOptions := cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}
	if s.cfg.Auth != nil && len(s.cfg.Auth.AllowedOrigins) > 0 {
		corsOptions.AllowedOrigins = s.cfg.Auth.AllowedOrigins
	}
	c := cors.New(corsOptions)
	r.Use(c.Handler)

	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "http api is up")
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		err := s.auth.Authenticate(r)
		if iferr.Respond(w, err) {
			s.RequestChannel <- domain.NewRequestError(err)

			return
		}

		output, err := output.NewOutput(r)
		if iferr.Respond(w, err) {
			s.RequestChannel <- domain.NewRequestError(err)

			return
		}

		req := domain.NewRequest(output)

		res := s.responseListener(req)
		if iferr.Respond(w, res.Error) {
			return
		}

		render.Respond(w, r, res.Message)
	})

	server.Handler = r
	err := server.ListenAndServe()
	switch err {
	case nil:
		return
	case http.ErrServerClosed:
		return
	default:
		iferr.Exit(errors.Wrapf(err, "failed serving on port %s", server.Addr))

		return
	}
}
