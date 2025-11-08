package middleware

import (
	"net/http"

	"social-network/app/repository"
)

func AuthenticateUser(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", err
	}

	var userID string
	err = repository.Db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", cookie.Value).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}
