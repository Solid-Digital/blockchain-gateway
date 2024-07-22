package apperr

import (
	"bytes"
	stderr "errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/unchainio/pkg/errors"

	"github.com/unchainio/interfaces/logger"
)

const defaultStatus = http.StatusInternalServerError

func New() *Error {
	return &Error{
		Status:      int64(defaultStatus),
		StatusText:  http.StatusText(defaultStatus),
		Message:     "",
		NamedErrors: nil,
		RequestID:   "",
		cause:       nil,
	}
}

type Error struct {
	// status
	Status     int64  `json:"status,omitempty"`
	StatusText string `json:"statusText,omitempty"`
	Code       int64  `json:"code,omitempty"`

	// message
	Message string `json:"message,omitempty"`

	// named errors
	NamedErrors map[string][]string `json:"namedErrors,omitempty"`
	// request ID
	RequestID string `json:"requestID,omitempty"`

	cause error `json:"-"`
}

func (e *Error) Is(err error) bool {
	e2, ok := err.(*Error)
	if !ok {
		return false
	}

	return e.Status == e2.Status && e.StatusText == e2.StatusText && e.Message == e2.Message
}

func copyMap(dst, src *map[string][]string) {
	*dst = make(map[string][]string, len(*src))
	for key, val := range *src {
		var outVal []string
		if val == nil {
			(*dst)[key] = nil
		} else {
			in, out := &val, &outVal
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
		(*dst)[key] = outVal
	}
}

func (e *Error) Copy() *Error {
	var namedErrors map[string][]string

	copyMap(&namedErrors, &e.NamedErrors)

	return &Error{
		NamedErrors: namedErrors,
		Message:     e.Message,
		RequestID:   e.RequestID,
		Status:      e.Status,
		StatusText:  e.StatusText,
		cause:       e.cause,
	}
}

func (e *Error) Wrap(err error) *Error {
	e2 := e.Copy()
	e2.cause = err
	return e2
}

func (e *Error) WithMessage(message string) *Error {
	e2 := e.Wrap(e)
	e2.Message = message
	return e2
}

func (e *Error) WithMessagef(message string, args ...interface{}) *Error {
	return e.WithMessage(fmt.Sprintf(message, args...))
}

func (e *Error) WithStatus(status int64) *Error {
	e.Status = status
	e.StatusText = http.StatusText(int(e.Status))
	return e
}

func (e *Error) WithCode(code int64) *Error {
	e.Code = code
	return e
}

func (e *Error) WithRequestID(requestID string) *Error {
	e.RequestID = requestID
	return e
}

func (e *Error) AddNamedErrors(name string, errs ...string) *Error {
	if e.NamedErrors == nil {
		e.NamedErrors = make(map[string][]string)
	}

	e.NamedErrors[name] = append(e.NamedErrors[name], errs...)

	return e
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		io.WriteString(s, e.Error())

		if s.Flag('+') {
			var stackTrace stackTracer
			if stderr.As(e, &stackTrace) {
				io.WriteString(s, fmt.Sprintf("%+v", stackTrace.StackTrace()))
			}
		}
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

// WriteResponse to the client
func (e *Error) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.WriteHeader(int(e.Status))

	if err := producer.Produce(rw, e); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

func (e *Error) Error() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%s", e.Message)

	if len(e.NamedErrors) != 0 {
		for name, v := range e.NamedErrors {
			fmt.Fprintf(buf, "\n* %s:\n", name)

			for _, vv := range v {
				fmt.Fprintf(buf, "* * %s\n", vv)
			}
		}
	}

	if e.cause != nil {
		fmt.Fprintf(buf, ": %s:\n", e.cause.Error())
	}

	return buf.String()
}

func (e *Error) Unwrap() error {
	return e.cause
}

func HandleError(err *Error, requestID string, log logger.Logger) bool {
	if err == nil {
		return false
	}

	err = err.WithRequestID(requestID)

	if err.Status == http.StatusInternalServerError {
		log.Errorf("%+v", err)
	} else {
		log.Debugf("%+v", err)
	}

	return true
}
