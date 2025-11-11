package model

import (
	"time"

	"social-network/pkg/db/sqlite"
)

func IsUserInGroup(userID, groupID string) (bool, error) {
	var isMember bool
	err := sqlite.Db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM group_members gm WHERE gm.user_id = ? AND gm.group_id = ?
		)
	`, userID, groupID).Scan(&isMember)
	return isMember, err
}

func DoesPostExistInGroup(postID, groupID string) (bool, error) {
	var exists bool
	err := sqlite.Db.QueryRow(`SELECT EXISTS(SELECT 1 FROM group_posts WHERE id = ? AND group_id = ?)`, postID, groupID).Scan(&exists)
	return exists, err
}

func GetExistingLikeGroup(userID, postID string) (string, error) {
	var likeID string
	err := sqlite.Db.QueryRow(`
		SELECT id FROM likesgroups
		WHERE user_id = ? AND liked_item_id = ? AND liked_item_type = 'post'
	`, userID, postID).Scan(&likeID)
	return likeID, err
}

func RemoveLikeGroup(likeID string) error {
	_, err := sqlite.Db.Exec(`DELETE FROM likesgroups WHERE id = ?`, likeID)
	return err
}

func AddLikeGroup(likeID, userID, postID string, createdAt time.Time) error {
	_, err := sqlite.Db.Exec(`
		INSERT INTO likesgroups (id, user_id, liked_item_id, liked_item_type, created_at)
		VALUES (?, ?, ?, 'post', ?)
	`, likeID, userID, postID, createdAt)
	return err
}
