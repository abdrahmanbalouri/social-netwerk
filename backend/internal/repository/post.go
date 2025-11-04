package repository

import (
	"database/sql"
	"strings"
	"time"
)

type Post struct {
	ID           string
	UserID       string
	Title        string
	Content      string
	ImagePath    string
	Visibility   string
	AllowedUsers []string
	CreatedAt    time.Time
}

// InsertPost inserts a post into the database
func InsertPost(db *sql.DB, post Post) error {
	_, err := db.Exec(`
		INSERT INTO posts (id, user_id, title, content, image_path, visibility, canseperivite)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		post.ID, post.UserID, post.Title, post.Content, post.ImagePath, post.Visibility,
		func() string {
			if post.Visibility == "private" {
				return strings.Join(post.AllowedUsers, ",")
			}
			return ""
		}(),
	)
	return err
}

// InsertAllowedUsers inserts allowed users for private posts
func InsertAllowedUsers(db *sql.DB, postID, authorID string, allowedUsers []string) error {
	for _, uid := range allowedUsers {
		if uid == "" {
			continue
		}
		var exists int
		err := db.QueryRow(`SELECT 1 FROM users WHERE id=?`, uid).Scan(&exists)
		if err != nil {
			continue
		}
		_, err = db.Exec(`
			INSERT INTO allowed_followers (user_id, post_id, allowed_user_id)
			VALUES (?, ?, ?)`, authorID, postID, uid)
		if err != nil {
			continue
		}
	}
	return nil
}
