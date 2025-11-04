package repository

import (
	"database/sql"
	"time"
)

type CommentRepo struct {
	ID        string
	Content   string
	CreatedAt time.Time
	FirstName string
	LastName  string
	MediaPath string
}

func FetchCommentsGroup(db *sql.DB, postID string, offset int) ([]CommentRepo, error) {
	query := `
	SELECT cg.id, cg.content, cg.created_at, u.first_name, u.last_name, cg.media_path
	FROM comments_groups cg
	JOIN users u ON cg.user_id = u.id
	WHERE cg.post_id = ?
	ORDER BY cg.created_at DESC
	LIMIT 10 OFFSET ?`

	rows, err := db.Query(query, postID, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []CommentRepo
	for rows.Next() {
		var c CommentRepo
		if err := rows.Scan(&c.ID, &c.Content, &c.CreatedAt, &c.FirstName, &c.LastName, &c.MediaPath); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}
func GetCommentByIDlast(db *sql.DB, commentId string) (*CommentRepo, error) {
    row := db.QueryRow(`
        SELECT 
            c.id, c.content, c.created_at,
            u.first_name, u.last_name, c.media_path
        FROM comments_groups c
        JOIN users u ON c.user_id = u.id
        WHERE c.id = ?
    `, commentId)

    var comment CommentRepo
    err := row.Scan(&comment.ID, &comment.Content, &comment.CreatedAt,
        &comment.FirstName, &comment.LastName, &comment.MediaPath)
    if err != nil {
        return nil, err
    }

    return &comment, nil
}