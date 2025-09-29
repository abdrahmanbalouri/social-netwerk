package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/internal/helper"
	"social-network/internal/repository"
	"strings"

	"github.com/google/uuid"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		fmt.Println(err)
 		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	type CommentRequest struct {
		PostID  string `json:"postId"`
		Content string `json:"content"`
	}

	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.PostID == "" || strings.TrimSpace(req.Content) == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}
	if len(req.Content) < 3 || len(req.Content) > 30 {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}
	sanitizedContent := helper.Skip(req.Content)
	row := repository.Db.QueryRow(`SELECT content FROM posts WHERE id = ?`, req.PostID)
	var exists string
	err1 := row.Scan(&exists)
	if err1 == sql.ErrNoRows {
		helper.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	} else if err1 != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	commentID := uuid.New().String()
	_, err = repository.Db.Exec(`
        INSERT INTO comments (id, post_id, user_id, content)
        VALUES (?, ?, ?, ?)`,
		commentID, req.PostID, userID, sanitizedContent)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"message":    "Comment created successfully",
		"comment_id": commentID,
	})
}
