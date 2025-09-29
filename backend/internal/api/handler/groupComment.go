package handlers

import (
	"database/sql"
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
	//get the grp's id
	query := `SELECT group_id FROM posts WHERE post_id = ?`
	var grpID string
	err = repository.Db.QueryRow(query, comment.PostID).Scan(&grpID)
	if err != nil{
		fmt.Println("Failed to get the group's id")
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get the group's id")
		return
	}
	// check if the user is a member of the group
	var isMember bool
	query = `SELECT EXISTS (SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?)`
	err = repository.Db.QueryRow(query, userID, grpID).Scan(&isMember)
	if err != nil {
		fmt.Println("Failed to check group membership")
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to check group membership")
		return
	}
	if !isMember {
		fmt.Println("The user is not a member of the group")
		helper.RespondWithError(w, http.StatusUnauthorized, "You are not a member of this group")
		return
	}

	if comment.PostID == "" || strings.TrimSpace(comment.Content) == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "PostID or commentComment is empty")
		return
	}
	if len(comment.Content) < 3 || len(comment.Content) > 30 {
		helper.RespondWithError(w, http.StatusBadRequest, "Comment's content is too short or too long")
		return
	}
	sanitizedContent := helper.Skip(comment.Content)
	var exists string
	err1 := repository.Db.QueryRow(`SELECT content FROM posts WHERE id = ?`, comment.PostID).Scan(&exists)
	if err1 == sql.ErrNoRows {
		helper.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	} else if err1 != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

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
