package casbin

import (
	"net/http"

	"bitbucket.org/unchain/ares/gen/dto"
	openapierrors "github.com/go-openapi/errors"
	"github.com/unchainio/pkg/errors"
)

func (e *Enforcer) Authorize(r *http.Request, u interface{}) error {
	err := e.authorize(r, u)

	if err != nil {
		e.log.Errorf("%+v", err)

		return openapierrors.New(http.StatusForbidden, errors.Message(err, 1))
	}

	return nil
}

func (e *Enforcer) authorize(r *http.Request, u interface{}) error {
	user := u.(*dto.User)
	ok := e.Enforce(user.ID, r.URL.Path, r.Method)

	if !ok {
		return errors.New("not authorized")
	}

	return nil
}
