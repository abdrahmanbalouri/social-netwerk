package post

import (
	"database/sql"

	"social-network/internal/repository"
)

func GetAllPosts(id string) ([]map[string]interface{}, error) {
	var rows *sql.Rows
	// Modified SQL query with JOIN to get nickname from users table and comment count from comments table
	if id != "" {
		rows, _ = repository.Db.Query(`
		SELECT 
			p.id, 
			p.user_id, 
			p.title, 
			p.content, 
			p.image_path, 
			p.created_at, 
			u.nickname,
			u.image AS profile,
			COUNT(c.id) AS comments_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE p.user_id = ?
		GROUP BY p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, u.nickname, u.image
		ORDER BY p.created_at DESC;
	`, id)
	} else {
		rows, _ = repository.Db.Query(`
		SELECT 
			p.id, 
			p.user_id, 
			p.title, 
			p.content, 
			p.image_path, 
			p.created_at, 
			u.nickname,
			u.image AS profile,
			COUNT(c.id) AS comments_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		GROUP BY p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, u.nickname, u.image
		ORDER BY p.created_at DESC;
	`)
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id string
		var userID string
		var title, content, imagePath, nickname, profile, createdAt string
		var commentsCount int

		// Scan the result into variables
		err := rows.Scan(&id, &userID, &title, &content, &imagePath, &createdAt, &nickname, &profile, &commentsCount)
		if err != nil {
			return nil, err
		}

		post := map[string]interface{}{
			"id":             id,
			"user_id":        userID,
			"title":          title,
			"content":        content,
			"image_path":     imagePath,
			"created_at":     createdAt,
			"author":         nickname,
			"profile":        profile,
			"comments_count": commentsCount,
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return the list of posts with nickname and comment count included
	return posts, nil
}
