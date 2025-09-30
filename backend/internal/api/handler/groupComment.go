package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

type CommentRequest struct {
	PostID  string `json:"postId"`
	Content string `json:"content"`
}

func CreateCommentGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		fmt.Println(err)
		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	var comment CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}
	// get the grp's id
	query := `SELECT p.group_id, EXISTS(SELECT 1 FROM group_members gm WHERE gm.user_id = ? AND gm.group_id = p.group_id) FROM posts p WHERE p.id = ?;`
	var grpID string
	var isMember bool
	err = repository.Db.QueryRow(query,userID, comment.PostID).Scan(&grpID, &isMember)
	if err != nil {
		fmt.Println("Failed to get the group's id or post not exist :", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get the group's id or post not exist")
		return
	}
	if !isMember{
		fmt.Println("User is not a member of the group")
		helper.RespondWithError(w, http.StatusUnauthorized, "User is not a member of the group")
		return
	}

	if comment.PostID == "" || strings.TrimSpace(comment.Content) == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "PostID or commentComment is empty")
		return
	}
	if len(comment.Content) < 2 || len(comment.Content) > 30 {
		helper.RespondWithError(w, http.StatusBadRequest, "Comment's content is too short or too long")
		return
	}
	sanitizedContent := helper.Skip(comment.Content)

	commentID := helper.GenerateUUID()
	_, err = repository.Db.Exec(`
        INSERT INTO comments (id, post_id, user_id, content) VALUES (?, ?, ?, ?)`,
		commentID, comment.PostID, userID, sanitizedContent)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "Comment created successfully for groups",
	})
}

func GetCommentGroup(w http.ResponseWriter, r *http.Request) {
}
