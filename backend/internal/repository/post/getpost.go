package post

import (
	"fmt"

	"social-network/internal/repository"
)

func GetAllPosts() ([]map[string]interface{}, error) {
	// Modified SQL query with JOIN to get nickname from the users table
	rows, err := repository.Db.Query(`
		SELECT 
			p.id, 
			p.user_id, 
			p.title, 
			p.content, 
			p.image_path, 
			p.created_at, 
			u.nickname 
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id string
		var userID int
		var title, content, image_path, nickname string
		var createdAt string

		// Scan the result into variables
		err := rows.Scan(&id, &userID, &title, &content, &image_path, &createdAt, &nickname)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}

		// Create a map for each post, including the nickname
		post := map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"title":      title,
			"content":    content,
			"image_path": image_path,
			"created_at": createdAt,
			"author":   nickname, 
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return the list of posts with the nickname included
	return posts, nil
}
