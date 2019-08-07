package auth

import (
	"errors"
	"fmt"
	"net/http"

	"git.sr.ht/~izzy/vbooks/apperr"
	"git.sr.ht/~izzy/vbooks/db"
)

// A Role represents a user's maximum permission level
type Role int

// Permissions work as follows:
// - Invalid:   Deactivated or otherwise invalid users
// - User:      Normal users
// - SuperUser: Has all server permissions
const (
	RoleInvalid Role = iota
	RoleUser
	RoleSuperUser
)

func (r Role) String() string {
	switch r {
	case RoleSuperUser:
		return "administrator"
	case RoleUser:
		return "user"
	default:
		return "invalid"
	}
}

// An Authorization represents a users' current authorization and authentication level
type Authorization struct {
	Username string
	Role     Role
}

// Get returns a user's Authorization based on a session and a minimum required Role.
// If the user
func Get(r *http.Request, minimumRole Role, db db.DB) (*Authorization, *apperr.Error) {
	user, err := verifySessionCookie(r, db)
	if err != nil {
		return nil, err
	}

	role := RoleUser
	if role < minimumRole {
		reason := fmt.Sprintf("User has role %s but required role is %s", role, minimumRole)
		return nil, apperr.New(errors.New("User '"+user+"Permission denied: "+reason),
			"permission denied", http.StatusUnauthorized)
	}

	return &Authorization{
		Username: user,
		Role:     RoleUser,
	}, nil

}
