package auth

import (
	"log"
	"net/http"

	"git.sr.ht/~izzy/vbooks/apperr"
)

type Handler func(user string, w http.ResponseWriter, r *http.Request) *apperr.Error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := verifySessionCookie(r)
	if err != nil {
		log.Println("Invalid auth attempt:", err.Wrapped)
		w.WriteHeader(err.StatusCode)
		return
	}

	if err := h(user, w, r); err != nil {
		log.Println(err)
		w.WriteHeader(err.StatusCode)
	}
}
