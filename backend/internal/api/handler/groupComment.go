package handlers

import (
	"database/sql"
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

func GetCommentGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {

		helper.RespondWithError(w, http.StatusNotFound, "post not found")
		return
	}
	PostID := parts[3]

	userID, err := helper.AuthenticateUser(r)
	if err != nil {

		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication failed")
		return
	}

	// Check for user's membership
	query := `SELECT p.group_id, EXISTS(SELECT 1 FROM group_members gm WHERE gm.user_id = ? AND gm.group_id = p.group_id) FROM group_posts p WHERE p.id = ?;`
	var grpID string
	var isMember bool
	err = repository.Db.QueryRow(query, userID, PostID).Scan(&grpID, &isMember)
	if err == sql.ErrNoRows {
		return
	} else if err != nil {
		return
	}
	if !isMember {

		helper.RespondWithError(w, http.StatusUnauthorized, "User is not a member of the group")
		return
	}

	// Fetch all the posts of this group
	query = `SELECT id, content, created_at FROM comments WHERE post_id = ?`
	rows, err := repository.Db.Query(query, PostID)
	if err != nil {

		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get comments")
		return
	}
	defer rows.Close()

	var CommentJson []Comment
	for rows.Next() {
		var c Comment
		err := rows.Scan(&c.ID, &c.Content, &c.CreatedAt)
		if err != nil {

			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to comments comments")
			return
		}
		CommentJson = append(CommentJson, c)
	}

	// Return the posts as a JSON response
	helper.RespondWithJSON(w, http.StatusOK, CommentJson)
}
