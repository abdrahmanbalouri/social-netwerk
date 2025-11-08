package model

import "social-network/app/repository"

func FollowRequestExists(userID, followerID string) bool {
	var exists bool
	err := repository.Db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM follow_requests 
			WHERE user_id = ? AND follower_id = ?
		)
	`, userID, followerID).Scan(&exists)
	return err == nil && exists
}

func AddFollower(userID, followerID string) error {
	_, err := repository.Db.Exec(`
		INSERT INTO followers (user_id, follower_id) 
		VALUES (?, ?)
	`, userID, followerID)
	return err
}

func DeleteFollowRequest(userID, followerID string) error {
	_, err := repository.Db.Exec(`
		DELETE FROM follow_requests 
		WHERE user_id = ? AND follower_id = ?
	`, userID, followerID)
	return err
}
