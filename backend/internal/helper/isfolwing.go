package helper

import (
	"net/http"

	"social-network/internal/repository"
)

func IsFollowing(r *http.Request, targetUserID string) (bool, error) {
	authUserID, err := AuthenticateUser(r)
	if err != nil {
		return false, err
	}
	var exists bool
	err = repository.Db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM followers 
			WHERE user_id = ? AND follower_id = ?
		)`, targetUserID, authUserID).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
