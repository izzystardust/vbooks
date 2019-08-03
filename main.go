package main

import (
	"fmt"
	"log"
	"net/http"

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

func main() {
	http.Handle("/login", AppHandler(handleLogin))
	http.Handle("/welcome", auth.Handler(welcomeHandler))

	log.Fatal(http.ListenAndServe(":3001", nil))
}

func welcomeHandler(user string, w http.ResponseWriter, r *http.Request) *apperr.Error {
	w.Write([]byte(fmt.Sprintf("Welcome, %s", user)))

	return nil
}
