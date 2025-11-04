package repository

import (
	"database/sql"
	"time"
)

type Like struct {
	ID           string
	UserID       string
	LikedItemID  string
	LikedItemType string
	CreatedAt    time.Time
}

// Check if user already liked the post
func UserLikedPost(db *sql.DB, userID, postID string) (bool, string, error) {
	var id string
	err := db.QueryRow(`
		SELECT id FROM likes 
		WHERE user_id = ? AND liked_item_id = ? AND liked_item_type = 'post'
	`, userID, postID).Scan(&id)

	if err == sql.ErrNoRows {
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}
	return true, id, nil
}

// Insert a new like
func InsertLike(db *sql.DB, like Like) error {
	_, err := db.Exec(`
		INSERT INTO likes (id, user_id, liked_item_id, liked_item_type, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, like.ID, like.UserID, like.LikedItemID, like.LikedItemType, like.CreatedAt)
	return err
}

// Remove like by ID
func DeleteLike(db *sql.DB, likeID string) error {
	_, err := db.Exec(`DELETE FROM likes WHERE id = ?`, likeID)
	return err
}
