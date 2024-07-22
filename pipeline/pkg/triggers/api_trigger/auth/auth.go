package auth

import "net/http"

type Service interface {
	Authenticate(*http.Request) error
}
