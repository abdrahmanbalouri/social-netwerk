package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func GetGroupPostByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	currentUserID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}
	postID := parts[3]
	fmt.Println(postID)
	query := `
	SELECT 
		gp.id, 
		gp.user_id, 
		gp.title, 
		gp.content, 
		gp.image_path, 
		gp.created_at, 
		u.first_name,
		u.last_name,
		u.image AS profile,
		COUNT(DISTINCT l.id) AS like_count,
		COUNT(DISTINCT CASE WHEN l.user_id = ? THEN l.id END) AS liked_by_user,
		COUNT(DISTINCT c.id) AS comments_count
	FROM group_posts gp
	JOIN users u ON gp.user_id = u.id
	LEFT JOIN likes l ON gp.id = l.liked_item_id AND l.liked_item_type = 'post'
	LEFT JOIN comments c ON gp.id = c.post_id
	WHERE gp.id = ?
	GROUP BY 
		gp.id, gp.user_id, gp.title, gp.content, gp.image_path, gp.created_at, 
		u.first_name, u.last_name, u.image
	`

	var (
		id            string
		userID        string
		title         string
		content       string
		imagePath     sql.NullString
		createdAt     time.Time
		firstName     string
		lastName      string
		profile       sql.NullString
		likeCount     int
		likedByUser   int
		commentsCount int
	)

	err = repository.Db.QueryRow(query, currentUserID, postID).Scan(
		&id, &userID, &title, &content, &imagePath, &createdAt,
		&firstName, &lastName, &profile,
		&likeCount, &likedByUser, &commentsCount,
	)

	if err == sql.ErrNoRows {
		helper.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
		return
	}

	post := map[string]interface{}{
		"id":             id,
		"user_id":        userID,
		"title":          title,
		"content":        content,
		"image_path":     nullStringToString(imagePath),
		"created_at":     createdAt.Format(time.RFC3339),
		"first_name":     firstName,
		"last_name":      lastName,
		"profile":        nullStringToString(profile),
		"like":           likeCount,
		"liked_by_user":  likedByUser > 0,
		"comments_count": commentsCount,
	}

	helper.RespondWithJSON(w, http.StatusOK, post)
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
