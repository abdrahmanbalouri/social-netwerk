package post

import (
	"fmt"

	"social-network/internal/repository"
)

func GetAllPosts() ([]map[string]interface{}, error) {
	// Modified SQL query with JOIN to get nickname from users table and comment count from comments table
	rows, err := repository.Db.Query(`
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
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id string
		var userID int
		var title, content, imagePath, nickname, profile, createdAt string
		var commentsCount int

		// Scan the result into variables
		err := rows.Scan(&id, &userID, &title, &content, &imagePath, &createdAt, &nickname, &profile, &commentsCount)
		if err != nil {
			fmt.Println("Error scanning row:", err)
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

	if err = rows.Err(); err != nil {
		fmt.Println("Error iterating rows:", err)
		return nil, err
	}

	// Return the list of posts with nickname and comment count included
	return posts, nil
}