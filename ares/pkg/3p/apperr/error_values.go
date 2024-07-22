package apperr

import (
	"net/http"
)

var Conflict = New().WithStatus(http.StatusConflict)
var NotFound = New().WithStatus(http.StatusNotFound)
var BadRequest = New().WithStatus(http.StatusBadRequest)
var Unauthorized = New().WithStatus(http.StatusUnauthorized)
var Forbidden = New().WithStatus(http.StatusForbidden)
var Internal = New().WithStatus(http.StatusInternalServerError)
