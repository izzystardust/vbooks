package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"git.sr.ht/~izzy/vbooks/apperr"
	"git.sr.ht/~izzy/vbooks/auth"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) *apperr.Error

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Println(err)
		w.WriteHeader(err.StatusCode)
	}
}

type Router struct {
}

func Start(addr string) error {

	http.Handle("/login", AppHandler(handleLogin))
	http.Handle("/welcome", auth.Handler(welcomeHandler))

	return http.ListenAndServe(addr, nil)
}

func welcomeHandler(user string, w http.ResponseWriter, r *http.Request) *apperr.Error {
	w.Write([]byte(fmt.Sprintf("Welcome, %s", user)))

	return nil
}

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
