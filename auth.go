package main

import (
	"encoding/json"
	"net/http"
	"time"

	"git.sr.ht/~izzy/vbooks/apperr"
	"git.sr.ht/~izzy/vbooks/auth"
)

func handleLogin(w http.ResponseWriter, r *http.Request) *apperr.Error {
	var creds auth.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		return &apperr.Error{err, http.StatusUnauthorized}
	}

	session, err := auth.NewSession(creds)
	if err != nil {
		if httperr, ok := err.(*apperr.Error); ok {
			return httperr
		} else {
			return apperr.InternalServerError(err)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:    auth.SessionCookieName,
		Value:   session,
		Expires: time.Now().Add(auth.ExpirationTime),
	})

	return nil
}
