package auth

import (
	"errors"
	"net/http"
	"vbooks/apperr"
)

var (
	Unauthorized = &apperr.Error{
		errors.New("Unauthorized"),
		http.StatusUnauthorized,
	}
)
