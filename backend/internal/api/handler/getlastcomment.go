package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func Getlastcommnet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusNotFound, "Comment not found")
		return
	}

	commentId := parts[3]

	row := repository.Db.QueryRow(`
		SELECT 
			c.id, 
			c.post_id, 
			c.user_id, 
			c.content, 
			c.created_at,
			u.nickname,
			u.image AS profile,
			c.media_path
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = ?
	`, commentId)

	var comment struct {
		ID        string `json:"id"`
		PostID    string `json:"post_id"`
		UserID    string `json:"user_id"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
		Author    string `json:"author"`
		Profile   string `json:"profile"`
		MediaPath string `json:"media_path,omitempty"`
	}

	err := row.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.Author, &comment.Profile, &comment.MediaPath)
	if err != nil {
		if err == sql.ErrNoRows {
			helper.RespondWithError(w, http.StatusNotFound, "Comment not found")
			return
		}
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(comment.MediaPath)

	helper.RespondWithJSON(w, http.StatusOK, comment)
}
