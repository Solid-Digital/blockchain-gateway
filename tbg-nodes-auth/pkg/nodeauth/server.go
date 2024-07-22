package nodeauth

import (
	"github.com/go-chi/chi"
	"github.com/unchain/tbg-nodes-auth/pkg/3p/sql"
	"github.com/unchainio/interfaces/logger"
)

// Server is the ares server struct
type Server struct {
	Meta *Metadata
	DB   *sql.DB
	Log  logger.Logger
}

// NewServer constructs a new Ares server
func NewServer(db *sql.DB, metadata *Metadata, log logger.Logger) (*Server, error) {
	server := &Server{
		DB:   db,
		Log:  log,
		Meta: metadata,
	}

	return server, nil
}

func (s *Server) Handler() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/permissioned/{uuid}", s.HandleAuthPermissioned)
	r.Get("/*", s.HandleAuth)

	return r
}
