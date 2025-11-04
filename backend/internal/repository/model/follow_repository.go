package model

import "social-network/internal/repository"

func GetUserPrivacy(id string) (string, error) {
	var p string
	err := repository.Db.QueryRow(`SELECT privacy FROM users WHERE id = ?`, id).Scan(&p)
	return p, err
}

func IsFollowing(userID, followerID string) (bool, error) {
	var c int
	err := repository.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE user_id = ? AND follower_id = ?`, userID, followerID).Scan(&c)
	return c > 0, err
}

func Follow(userID, followerID string) error {
	_, err := repository.Db.Exec(`INSERT INTO followers (user_id, follower_id) VALUES (?, ?)`, userID, followerID)
	return err
}

func Unfollow(userID, followerID string) error {
	_, err := repository.Db.Exec(`DELETE FROM followers WHERE user_id = ? AND follower_id = ?`, userID, followerID)
	return err
}

func IsPending(userID, followerID string) (bool, error) {
	var c int
	err := repository.Db.QueryRow(`SELECT COUNT(*) FROM follow_requests WHERE user_id = ? AND follower_id = ?`, userID, followerID).Scan(&c)
	return c > 0, err
}

func CreateFollowRequest(userID, followerID string) error {
	_, err := repository.Db.Exec(`INSERT INTO follow_requests (user_id, follower_id) VALUES (?, ?)`, userID, followerID)
	return err
}

func CancelFollowRequest(userID, followerID string) error {
	_, err := repository.Db.Exec(`DELETE FROM follow_requests WHERE user_id = ? AND follower_id = ?`, userID, followerID)
	return err
}

func CountFollowers(id string) int {
	var c int
	repository.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE user_id = ?`, id).Scan(&c)
	return c
}

func CountFollowing(id string) int {
	var c int
	repository.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE follower_id = ?`, id).Scan(&c)
	return c
}
