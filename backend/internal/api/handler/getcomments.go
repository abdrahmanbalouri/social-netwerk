package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	

	//"social-network/internal/database"
	"social-network/internal/helper"
	"social-network/internal/repository"
)

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		helper.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	postID := parts[3]
	offsetStr := parts[4]

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
	_, err = helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	rows, err := repository.Db.Query(`
        SELECT c.id, c.content, c.created_at, u.first_name , u.last_name, c.media_path
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = ?
        ORDER BY c.created_at DESC
        LIMIT 10 OFFSET ?`, postID, offset)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch comments")
		return
	}
	defer rows.Close()

	type Comment struct {
		ID        string    `json:"id"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`

		First_name    string    `json:"first_name"`
		Last_name    string    `json:"last_name"`
		MediaPath string    `json:"media_path,omitempty"`
	}

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.First_name,&comment.Last_name, &comment.MediaPath)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to process comments")
			return
		}
		comments = append(comments, comment)
	}

	helper.RespondWithJSON(w, http.StatusOK, comments)
}
