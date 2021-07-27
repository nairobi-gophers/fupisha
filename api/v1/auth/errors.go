package auth

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

//ErrInvalidEmailOrPassword a wrong email or password field.
var ErrInvalidEmailOrPassword = errors.New("invalid email or password")

//ErrUnknownLogin an unregistered email field.
var ErrUnknownLogin = errors.New("email not registered")

//ErrMissingToken a missing authorization header with the Bearer token.
var ErrMissingToken = errors.New("missing authorization header")

//ErrNoSuchAccount a non-existent user account.
var ErrNoSuchAccount = errors.New("no such account")

//ErrEmailTaken an already registered email.
var ErrEmailTaken = errors.New("that email is taken")

//ErrLoginToken an invalid or expired .
var ErrLoginToken = errors.New("invalid or expired login token")

//ErrAPIUnsupported an unsupported api version.
var ErrAPIUnsupported = errors.New("unsupported api version")

//ErrMissingAPIVersion a missing api version header with the version text.
var ErrMissingAPIVersion = errors.New("missing api version header")

//ErrInvalidVerificationToken an expired or invalid verification token.
var ErrInvalidVerificationToken = errors.New("invalid or expired verification token")

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render sets the application-specific error code in AppCode.
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest returns status 422 Unprocessable Entity including error message.
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     http.StatusText(http.StatusUnprocessableEntity),
		ErrorText:      err.Error(),
	}
}

// ErrUnsupportedAPIVersion returns status 400 Bad Request including error message.
func ErrUnsupportedAPIVersion(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     http.StatusText(http.StatusBadRequest),
		ErrorText:      err.Error(),
	}
}

// ErrDuplicateField returns status 409 Status Conflict including error message.
func ErrDuplicateField(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusConflict,
		StatusText:     http.StatusText(http.StatusConflict),
		ErrorText:      err.Error(),
	}
}

// ErrUnauthorized renders status 401 Unauthorized with custom error message.
func ErrUnauthorized(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     http.StatusText(http.StatusUnauthorized),
		ErrorText:      err.Error(),
	}
}

// The list of default error types without specific error message.
var (
	ErrInternalServerError = &ErrResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     http.StatusText(http.StatusInternalServerError),
	}
	ErrForbidden = &ErrResponse{
		HTTPStatusCode: http.StatusForbidden,
		StatusText:     http.StatusText(http.StatusForbidden),
	}
)
