package service

import (
	"errors"
	"time"

	"social-network/app/repository"
	"social-network/app/repository/model"
)

// ValidateSession checks if a session token is valid and not expired
func ValidateSession(token string) (string, error) {
	session, err := model.GetSession(repository.Db, token)
	if err != nil {
		return "", errors.New("unauthorized")
	}

	if session.ExpiresAt.Before(time.Now()) {
		model.DeleteSession(repository.Db, token)
		return "", errors.New("unauthorized")
	}

	return session.UserID, nil
}
