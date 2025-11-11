package model

import "social-network/pkg/db/sqlite"

func GetUserPrivacy(id string) (string, error) {
	var p string
	err := sqlite.Db.QueryRow(`SELECT privacy FROM users WHERE id = ?`, id).Scan(&p)
	return p, err
}

func IsFollowing(userID, followerID string) (bool, error) {
	var c int
	err := sqlite.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE user_id = ? AND follower_id = ?`, userID, followerID).Scan(&c)
	return c > 0, err
}

func Follow(userID, followerID string) error {
	_, err := sqlite.Db.Exec(`INSERT INTO followers (user_id, follower_id) VALUES (?, ?)`, userID, followerID)
	return err
}

func Unfollow(userID, followerID string) error {
	_, err := sqlite.Db.Exec(`DELETE FROM followers WHERE user_id = ? AND follower_id = ?`, userID, followerID)
	return err
}

func IsPending(userID, followerID string) (bool, error) {
	var c int
	err := sqlite.Db.QueryRow(`SELECT COUNT(*) FROM follow_requests WHERE user_id = ? AND follower_id = ?`, userID, followerID).Scan(&c)
	return c > 0, err
}

func CreateFollowRequest(userID, followerID string) error {
	_, err := sqlite.Db.Exec(`INSERT INTO follow_requests (user_id, follower_id) VALUES (?, ?)`, userID, followerID)
	return err
}

func CancelFollowRequest(userID, followerID string) error {
	_, err := sqlite.Db.Exec(`DELETE FROM follow_requests WHERE user_id = ? AND follower_id = ?`, userID, followerID)
	return err
}

func CountFollowers(id string) int {
	var c int
	sqlite.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE user_id = ?`, id).Scan(&c)
	return c
}

func CountFollowing(id string) int {
	var c int
	sqlite.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE follower_id = ?`, id).Scan(&c)
	return c
}
