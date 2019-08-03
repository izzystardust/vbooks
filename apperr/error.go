package apperr

import "net/http"

type Error struct {
	Wrapped    error // Internal
	StatusCode int   // End User facing
}

func (ae *Error) Error() string {
	return ae.Wrapped.Error()
}

func Unauthorized(err error) *Error {
	return &Error{err, http.StatusUnauthorized}
}

func InternalServerError(err error) *Error {
	return &Error{err, http.StatusInternalServerError}
}

func BadRequest(err error) *Error {
	return &Error{err, http.StatusBadRequest}
}
