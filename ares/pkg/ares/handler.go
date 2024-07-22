package ares

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-openapi/runtime"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"
	"github.com/go-openapi/errors"

	"github.com/unchainio/pkg/iferr"

	"bitbucket.org/unchain/ares/gen/api"
)

type Middleware func(next http.Handler) http.Handler

func (ares *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, err := ares.Handler()

	if iferr.Respond(w, err) {
		return
	}

	h.ServeHTTP(w, r)
}

// Handler returns the http Handler of the ares server
func (ares *Server) Handler() (http.Handler, error) {
	h, api, err := api.HandlerAPI(*ares.Handlers)

	if err != nil {
		return nil, err
	}
	api.ServeError = func(w http.ResponseWriter, r *http.Request, err error) {
		ServeError(w, r, err)
	}

	return ares.Middleware(h), nil
}

// ServeError the error handler interface implementation
func ServeError(rw http.ResponseWriter, r *http.Request, err error) {
	rw.Header().Set("Content-Type", "application/json")

	h, appErr := asAppErr(err)
	for k, v := range h {
		for _, vv := range v {
			rw.Header().Add(k, vv)
		}
	}

	if r == nil || r.Method != http.MethodHead {
		appErr.WriteResponse(rw, runtime.JSONProducer())
	} else {
		rw.WriteHeader(int(appErr.Status))
	}
}

func asAppErr(err error) (http.Header, *apperr.Error) {
	switch e := err.(type) {
	case *errors.CompositeError:
		er := flattenComposite(e)
		// strips composite errors to first element only
		if len(er.Errors) > 0 {
			return asAppErr(er.Errors[0])
		} else {
			// guard against empty CompositeError (invalid construct)
			return asAppErr(nil)
		}
	case *errors.MethodNotAllowedError:
		h := make(http.Header)
		h.Add("Allow", strings.Join(err.(*errors.MethodNotAllowedError).Allowed, ","))

		return h, apperr.New().Wrap(err).
			WithCode(int64(e.Code())).
			WithStatus(int64(asHTTPCode(int(e.Code()))))
	case *errors.Validation:
		return nil, apperr.New().
			Wrap(err).
			WithCode(int64(e.Code())).
			WithStatus(int64(asHTTPCode(int(e.Code())))).
			AddNamedErrors(e.Name, e.Error())
	case errors.Error:
		value := reflect.ValueOf(e)
		if value.Kind() == reflect.Ptr && value.IsNil() {
			return nil, apperr.New().
				Wrap(err).
				WithCode(http.StatusInternalServerError).
				WithStatus(http.StatusInternalServerError).
				WithMessage("Unknown error")
		}

		return nil, apperr.New().
			Wrap(err).
			WithCode(int64(e.Code())).
			WithStatus(int64(asHTTPCode(int(e.Code())))).
			WithMessage(e.Error())
	case nil:
		return nil, apperr.New().
			WithCode(http.StatusInternalServerError).
			WithStatus(http.StatusInternalServerError).
			WithMessage("Unknown error")
	default:
		return nil, apperr.New().
			Wrap(err).
			WithCode(http.StatusInternalServerError).
			WithStatus(http.StatusInternalServerError).
			WithMessage(err.Error())
	}
}

func flattenComposite(errs *errors.CompositeError) *errors.CompositeError {
	var res []error
	for _, er := range errs.Errors {
		switch e := er.(type) {
		case *errors.CompositeError:
			if len(e.Errors) > 0 {
				flat := flattenComposite(e)
				if len(flat.Errors) > 0 {
					res = append(res, flat.Errors...)
				}
			}
		default:
			if e != nil {
				res = append(res, e)
			}
		}
	}
	return errors.CompositeValidationError(res...)
}

func asHTTPCode(input int) int {
	if input >= 600 {
		return errors.DefaultHTTPCode
	}
	return input
}
