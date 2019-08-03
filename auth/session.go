package auth

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"git.sr.ht/~izzy/vbooks/apperr"

	uuid "github.com/satori/go.uuid"
)

const (
	// The SessionCookieName is the name of the session cookie expected
	// for authentication
	SessionCookieName = "vbooks_session"

	// ExpirationTime is the duration a session cookie is valid for once issued
	ExpirationTime = 10 * time.Minute
)

var (
	sessions    = map[string]session{}
	sessionLock sync.Mutex
)

// Credentials represents the JSON object used for password authentication
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//TODO: this goes into a database
//TODO: Yeah that's not how you store passwords
var users = map[string]string{
	"izzy": "test99**",
}

type session struct {
	user    string
	expires time.Time
}

// NewSession
func NewSession(user Credentials) (string, error) {
	if err := verifyLogin(user); err != nil {
		return "", err
	}
	sessionToken, err := uuid.NewV4()
	if err != nil {
		return "", &apperr.Error{err, http.StatusInternalServerError}
	}

	sessions[sessionToken.String()] = session{
		user.Username,
		time.Now().Add(ExpirationTime),
	}

	return sessionToken.String(), nil
}

func verifyLogin(user Credentials) error {
	if len(user.Username) == 0 || len(user.Password) == 0 {
		return apperr.BadRequest(errors.New("Missing or empty field in JSON"))
	}
	expectedPass, exists := users[user.Username]
	if !exists || expectedPass != user.Password {
		return apperr.Unauthorized(errors.New("Invalid login attempt for " + user.Username))
	}
	return nil
}

func (s session) isExpired() bool {
	return time.Now().After(s.expires)
}

func verifySessionCookie(r *http.Request) (string, *apperr.Error) {
	c, err := r.Cookie(SessionCookieName)
	if err == http.ErrNoCookie {
		return "", &apperr.Error{err, http.StatusUnauthorized}
	} else if err != nil {
		return "", &apperr.Error{err, http.StatusBadRequest}
	}

	//TODO: lock
	session, ok := sessions[c.Value]
	if !ok {
		return "", &apperr.Error{errors.New("Session does not exist"), http.StatusUnauthorized}
	}

	//TODO: remove expired sessions
	if session.isExpired() {
		return "", &apperr.Error{errors.New("Session expired"), http.StatusUnauthorized}
	}

	return session.user, nil
}
