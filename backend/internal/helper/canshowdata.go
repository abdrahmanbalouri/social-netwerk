package helper

import (
	"database/sql"
	"fmt"

	"social-network/internal/repository"
)

func Canshowdata(userID, postID string) (bool, error) {
	var (
		postUserID  string
		visibility  string
		userPrivacy string
	)
	err := repository.Db.QueryRow(`
		SELECT p.user_id, p.visibility, u.privacy
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`, postID).Scan(&postUserID, &visibility, &userPrivacy)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("post not found")
		}
		return false, err
	}

	if userID == postUserID {
		return true, nil
	}

	if visibility == "public" && userPrivacy == "public" {
		return true, nil
	}

	if visibility == "public" && userPrivacy == "private" {
		var exists int
		err := repository.Db.QueryRow(`
			SELECT 1 FROM followers 
			WHERE user_id = ? AND follower_id = ?
		`, userID, postUserID).Scan(&exists)
		if err == nil {
			return true, nil
		}
	}

	if visibility == "almost_private" {
		var exists int
		err := repository.Db.QueryRow(`
			SELECT 1 FROM followers 
			WHERE user_id = ? AND follower_id = ?
		`, userID, postUserID).Scan(&exists)
		if err == nil {
			return true, nil
		}
	}

	if visibility == "private" {
		var exists int
		err := repository.Db.QueryRow(`
			SELECT 1 FROM allowed_followers
			WHERE allowed_user_id = ? AND post_id = ? AND user_id = ?
		`, userID, postID, postUserID).Scan(&exists)
		if err == nil {
			return true, nil
		}
	}

	return false, nil
}
