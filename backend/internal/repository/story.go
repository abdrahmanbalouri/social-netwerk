package repository

import (
	"database/sql"
)

type Story struct {
	UserID   string
	Content  string
	ImageURL string
	BgColor  string
}

// InsertStory inserts a story into the database
func InsertStory(db *sql.DB, story Story) error {
	_, err := db.Exec(`
		INSERT INTO stories (user_id, content, image_url, bg_color)
		VALUES (?, ?, ?, ?)`,
		story.UserID, story.Content, story.ImageURL, story.BgColor,
	)
	return err
}
