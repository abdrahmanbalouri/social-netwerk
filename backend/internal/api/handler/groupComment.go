package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

type CommentRequest struct {
	PostID  string `json:"postId"`
	Content string `json:"content"`
}
type Comment struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	AuthorID  string    `json:"author_id"`
}

// type FetchComment struct{
// 	PostID string `json:"postId"`
// }

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
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format1")
		return
	}
	// get the grp's id
	query := `SELECT p.group_id, EXISTS(SELECT 1 FROM group_members gm WHERE gm.user_id = ? AND gm.group_id = p.group_id) FROM posts p WHERE p.id = ?;`
	var grpID string
	var isMember bool
	err = repository.Db.QueryRow(query, userID, comment.PostID).Scan(&grpID, &isMember)
	if err != nil {
		fmt.Println("Failed to get the group's id or post not exist :", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get the group's id or post not exist")
		return
	}
	if !isMember {
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
	fmt.Println("INSIDE GET COMMENTS")
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// var newComment FetchComment
	// if err := json.NewDecoder(r.Body).Decode(&newComment); err != nil {
	// 	fmt.Println("error is :", err)
	// 	helper.RespondWithError(w, http.StatusBadRequest, "Invalid request formattt")
	// 	return
	// }
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		fmt.Println("POST NOT FOOUND")
		helper.RespondWithError(w, http.StatusNotFound, "post not found")
		return
	}
	PostID := parts[3]

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		fmt.Println("Authentication failed : ", err)
		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication failed")
		return
	}

	// Check for user's membership
	query := `SELECT p.group_id, EXISTS(SELECT 1 FROM group_members gm WHERE gm.user_id = ? AND gm.group_id = p.group_id) FROM group_posts p WHERE p.id = ?;`
	var grpID string
	var isMember bool
	err = repository.Db.QueryRow(query, userID, PostID).Scan(&grpID, &isMember)
	if err == sql.ErrNoRows {
		fmt.Println("No post found with that ID")
		return
	} else if err != nil {
		fmt.Println("Database error:", err)
		return
	}
	if !isMember {
		fmt.Println("User is not a member of the group")
		helper.RespondWithError(w, http.StatusUnauthorized, "User is not a member of the group")
		return
	}

	// Fetch all the posts of this group
	query = `SELECT id, content, created_at FROM comments WHERE post_id = ?`
	rows, err := repository.Db.Query(query, PostID)
	if err != nil {
		fmt.Println("Failed to get comments : ", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get comments")
		return
	}
	defer rows.Close()

	var CommentJson []Comment
	for rows.Next() {
		var c Comment
		err := rows.Scan(&c.ID, &c.Content, &c.CreatedAt)
		if err != nil {
			fmt.Println("Failed to get comments : ", err)
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to comments comments")
			return
		}
		CommentJson = append(CommentJson, c)
	}

	// Return the posts as a JSON response
	helper.RespondWithJSON(w, http.StatusOK, CommentJson)
}
