package repository

import (
	"database/sql"
)

type Comment struct {
	ID        string
	PostID    string
	UserID    string
	Content   string
	MediaPath string
	GroupID   string
	IsGroup   bool
}

// InsertComment saves the comment in the correct table
func InsertComment(db *sql.DB, c Comment) error {
	if c.IsGroup {
		_, err := db.Exec(`
			INSERT INTO comments_groups (id, post_id, user_id, content, media_path)
			VALUES (?, ?, ?, ?, ?)`,
			c.ID, c.PostID, c.UserID, c.Content, c.MediaPath,
		)
		return err
	} else {
		_, err := db.Exec(`
			INSERT INTO comments (id, post_id, user_id, content, media_path)
			VALUES (?, ?, ?, ?, ?)`,
			c.ID, c.PostID, c.UserID, c.Content, c.MediaPath,
		)
		return err
	}
}
