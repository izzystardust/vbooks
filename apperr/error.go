package apperr

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// An Error represents an expected error in processing a user request.
// It includes information that is intended to be sent to a consumer of the API.
type Error struct {
	Wrapped    error  // Internal
	Reason     string // End user facing
	StatusCode int    // End User facing
	file       string
	line       int
}

// New returns an Error that includes wraps the given error with additional context
// If err is nil, New creates an error containing reason.
func New(err error, reason string, statuscode int) *Error {
	return new(err, reason, statuscode)
}

func new(err error, reason string, statuscode int) *Error {
	if err == nil {
		err = errors.New(reason)
	}
	_, file, line, ok := runtime.Caller(2)
	if ok {
		file = projectRelative(file)
	} else {
		file = "caller unavailiable"
	}

	return &Error{
		Wrapped:    err,
		Reason:     reason,
		StatusCode: statuscode,
		file:       file,
		line:       line,
	}
}

func projectRelative(absPath string) string {
	project := strings.Index(absPath, "vbooks")
	if project < 0 {
		return absPath
	}

	return absPath[project:]
}

func (ae *Error) Error() string {
	return fmt.Sprintf("%s:%d: %v (HTTP %d \"%s\")",
		ae.file, ae.line, ae.Wrapped, ae.StatusCode, ae.Reason)
}

// HTTPError writes error information to the header of w. If err is an apperr.Error,
// HTTPError writes the reason to the body as well as setting the Reason from the Error.
// If err is any other kind of error, HTTPError writes http.StatusInternalServerError.
func HTTPError(w http.ResponseWriter, err error) {
	if ae, ok := err.(*Error); ok {
		http.Error(w, ae.Reason, ae.StatusCode)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
