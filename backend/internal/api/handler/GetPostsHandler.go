package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("22222")
	if r.Method != "GET" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// Authenticate user
	UserId, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Extract post ID from URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}
	postID := parts[3]

	// Define struct for post response
	row := repository.Db.QueryRow(`
		SELECT 
			p.id, 
			p.user_id, 
			p.title, 
			p.content, 
			p.image_path, 
			p.created_at, 
			u.first_name, u.last_name,
			u.image AS profile,
			COUNT(DISTINCT l.id) AS like_count,
			COUNT(DISTINCT CASE WHEN l.user_id = ? THEN l.id END) AS liked_by_user,
			COUNT(DISTINCT c.id) AS comments_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN likes l ON p.id = l.liked_item_id AND l.liked_item_type = 'post'
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE p.id = ?
		GROUP BY p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, 	u.first_name, u.last_name, u.image
	`, UserId, postID)

	var id, userID, title, content, imagePath, first_name , last_name, profile, createdAt string
	var likeCount, likedByUser, commentsCount int

	// Scan the result into variables
	err = row.Scan(&id, &userID, &title, &content, &imagePath, &createdAt, &first_name,&last_name, &profile, &likeCount, &likedByUser, &commentsCount)
	if err != nil {
		if err == sql.ErrNoRows {
			helper.RespondWithError(w, http.StatusNotFound, "Post not found")
		}
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve post")
	}

	post := map[string]interface{}{
		"id":             id,
		"user_id":        userID,
		"title":          title,
		"content":        content,
		"image_path":     imagePath,
		"created_at":     createdAt,
		"first_name":         first_name,
		"last_name":         last_name,

		"profile":        profile,
		"like":           likeCount,
		"liked_by_user":  likedByUser > 0,
		"comments_count": commentsCount,
	}

	helper.RespondWithJSON(w, http.StatusOK, post)
}
