package service

import (
	"errors"
	"time"

	"social-network/internal/repository"
)

// ValidateSession checks if a session token is valid and not expired
func ValidateSession(token string) (string, error) {
	session, err := repository.GetSession(repository.Db, token)
	if err != nil {
		return "", errors.New("unauthorized")
	}

	if session.ExpiresAt.Before(time.Now()) {
		repository.DeleteSession(repository.Db, token)
		return "", errors.New("unauthorized")
	}

	return session.UserID, nil
}
