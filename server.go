package vbooks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"git.sr.ht/~izzy/vbooks/apperr"
	"git.sr.ht/~izzy/vbooks/auth"
	"git.sr.ht/~izzy/vbooks/db"
)

// A Config holds types used for dependency injection into handlers
type Config struct {
	DB db.DB
}

// AppHandler will be used to inject dependencies into the handlers, such as the database.
// If the handler returns an error, it shouldn't write to the http.ResponseWriter - the error
// will be used to set the response headers.
type AppHandler struct {
	Handler func(w http.ResponseWriter, r *http.Request, c Config) error
	Config  Config
}

func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //TODO: debug only
	if err := ah.Handler(w, r, ah.Config); err != nil {
		apperr.HTTPError(w, err)
		log.Println(err)
	}
}

// Start registers handlers and starts the application server
func Start(addr string) error {

	http.Handle("/login", AppHandler{handleLogin, Config{}})
	http.Handle("/welcome", AppHandler{welcomeHandler, Config{}})

	return http.ListenAndServe(addr, nil)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request, c Config) error {
	auth, err := auth.Get(r, auth.RoleUser, nil)
	if err != nil {
		return err
	}
	w.Write([]byte(fmt.Sprintf("Welcome, %s", auth.Username)))

	return nil
}

func handleLogin(w http.ResponseWriter, r *http.Request, c Config) error {
	var creds auth.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		return apperr.New(err, "", http.StatusBadRequest)
	}

	session, err := auth.NewSession(creds)
	if err != nil {
		if httperr, ok := err.(*apperr.Error); ok {
			return httperr
		}
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    auth.SessionCookieName,
		Value:   session,
		Expires: time.Now().Add(auth.ExpirationTime),
	})

	log.Println("User '" + creds.Username + "' logged in")
	return nil
}
