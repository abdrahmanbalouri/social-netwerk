package model

import (
	"database/sql"
	"fmt"
	"time"
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

type CommentAPI struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Profile   string    `json:"profile"`
	MediaPath string    `json:"media_path,omitempty"`
}

// InsertComment saves the comment in the correct table
func InsertComment(db *sql.DB, c Comment) error {
	if c.IsGroup {
		fmt.Println("hahahahahah")
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

// GetComments returns a list of comments for a post with offset
func GetComments(db *sql.DB, postID string, offset int) ([]CommentAPI, error) {
	query := `
        SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.first_name, u.last_name, u.image, c.media_path
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = ?
        ORDER BY c.created_at DESC
        LIMIT 10 OFFSET ?`
	rows, err := db.Query(query, postID, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []CommentAPI
	for rows.Next() {
		var c CommentAPI
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.FirstName, &c.LastName, &c.Profile, &c.MediaPath); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// GetCommentByID returns a single comment by ID
func GetCommentByID(db *sql.DB, commentID string) (CommentAPI, error) {
	var comment CommentAPI
	row := db.QueryRow(`
		SELECT 
			c.id, 
			c.post_id, 
			c.user_id, 
			c.content, 
			c.created_at,
			u.first_name,
			u.last_name,
			u.image AS profile,
			c.media_path
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = ?`, commentID)

	err := row.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content,
		&comment.CreatedAt, &comment.FirstName, &comment.LastName, &comment.Profile, &comment.MediaPath)
	return comment, err
}
