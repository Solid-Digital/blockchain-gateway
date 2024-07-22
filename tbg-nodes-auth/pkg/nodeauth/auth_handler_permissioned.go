package nodeauth

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/volatiletech/null"

	"github.com/unchain/tbg-nodes-auth/gen/orm"

	"github.com/unchainio/pkg/iferr"
)

func (s *Server) HandleAuthPermissioned(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	user, pass, ok := r.BasicAuth()

	if !ok {
		http.Error(w, "No basic auth was provided", http.StatusBadRequest)
		return
	}

	var foundCreds bool
	err := s.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		var err error

		foundCreds, err = orm.PermissionedEthereumNetworksBasicauthCreds(
			orm.PermissionedEthereumNetworksBasicauthCredWhere.ExternalInterfaceUUID.EQ(null.StringFrom(uuid)),
			orm.PermissionedEthereumNetworksBasicauthCredWhere.Username.EQ(null.StringFrom(user)),
			orm.PermissionedEthereumNetworksBasicauthCredWhere.Password.EQ(null.StringFrom(pass)),
		).Exists(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})
	if iferr.Respond(w, err, &iferr.ResponseOpts{Code: http.StatusInternalServerError}) {
		return
	}

	if !foundCreds {
		http.Error(w, "Request not allowed", http.StatusForbidden)
		return
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "Request allowed\n")

	s.Log.Debugf("Request allowed")
}
