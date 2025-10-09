package handlers

import (
	"fmt"
	"net/http"
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
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	postID := parts[3]
	Userid, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	can, err := helper.CanViewComments(Userid, postID)
	fmt.Println(can,"--------------------")
	if !can {
		helper.RespondWithError(w, http.StatusForbidden, "You do not have permission to view comments on this post")
		return
	}
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Error checking permissions")
		return
	}

	rows, err := repository.Db.Query(`
        SELECT c.id, c.content, c.created_at, u.nickname
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = ?
        ORDER BY c.created_at ASC`, postID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch comments")
		return
	}
	defer rows.Close()

	type Comment struct {
		ID        string    `json:"id"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
		Author    string    `json:"author"`
	}

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.Author)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to process comments")
			return
		}
		comments = append(comments, comment)
	}

	helper.RespondWithJSON(w, http.StatusOK, comments)
}
