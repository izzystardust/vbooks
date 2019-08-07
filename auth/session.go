package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"git.sr.ht/~izzy/vbooks/apperr"
	"git.sr.ht/~izzy/vbooks/db"

	uuid "github.com/satori/go.uuid"
)

const (
	// SessionCookieName is the name of the session cookie expected for authentication
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

type Token struct {
	Session string
	Expires time.Time
}

func (t Token) MarshalJSON() ([]byte, error) {
	basicToken := struct {
		SessionToken string `json:"sessiontoken"`
		Expires      string `json:"expires"`
	}{
		SessionToken: t.Session,
		Expires:      t.Expires.Format(time.RFC3339),
	}

	return json.Marshal(basicToken)
}

// NewSession creates a new session and inserts it into the session table
func NewSession(user Credentials) (*Token, error) {
	if err := verifyLogin(user); err != nil {
		return nil, err
	}
	sessionToken := uuid.NewV4().String()
	expires := time.Now().Add(ExpirationTime)

	sessions[sessionToken] = session{user.Username, expires}

	return &Token{sessionToken, expires}, nil
}

func verifyLogin(user Credentials) error {
	if len(user.Username) == 0 || len(user.Password) == 0 {
		return apperr.New(nil, "Missing or empty field in JSON", http.StatusBadRequest)
	}
	expectedPass, exists := users[user.Username]
	if !exists || expectedPass != user.Password {
		return apperr.New(errors.New("Invalid login attempt for "+user.Username),
			"Invalid username or password", http.StatusUnauthorized)
	}
	return nil
}

func (s session) isExpired() bool {
	return time.Now().After(s.expires)
}

func verifySessionCookie(r *http.Request, db db.DB) (string, *apperr.Error) {
	c, err := r.Cookie(SessionCookieName)
	if err == http.ErrNoCookie {
		return "", apperr.New(err, "Session cookie missing", http.StatusUnauthorized)
	} else if err != nil {
		return "", apperr.New(err, "", http.StatusBadRequest)
	}

	//TODO: lock
	session, ok := sessions[c.Value]
	if !ok {
		return "", apperr.New(nil, "Invalid session cookie", http.StatusUnauthorized)
	}

	//TODO: remove expired sessions
	if session.isExpired() {
		return "", apperr.New(nil, "Session expired", http.StatusUnauthorized)
	}

	return session.user, nil
}
