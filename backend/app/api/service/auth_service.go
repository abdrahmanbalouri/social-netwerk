package service

import (
	"errors"
	"time"

	"social-network/app/repository/model"
	"social-network/pkg/db/sqlite"
)

// ValidateSession checks if a session token is valid and not expired
func ValidateSession(token string) (string, error) {
	session, err := model.GetSession(sqlite.Db, token)
	if err != nil {
		return "", errors.New("unauthorized")
	}

	if session.ExpiresAt.Before(time.Now()) {
		model.DeleteSession(sqlite.Db, token)
		return "", errors.New("unauthorized")
	}

	return session.UserID, nil
}
