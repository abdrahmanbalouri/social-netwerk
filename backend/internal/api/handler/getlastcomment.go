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
		helper.RespondWithError(w, http.StatusNotFound, "coment not found ")
		return
	}

	coomentId := parts[3]
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	fmt.Println(userID)
	row := repository.Db.QueryRow(`
		SELECT 
			c.id, 
			c.post_id, 
			c.user_id, 
			c.content, 
			c.created_at,
			u.nickname,
			u.image AS profile
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = ?
	`, coomentId)

	var comment struct {
		ID        string `json:"id"`
		PostID    string `json:"post_id"`
		UserID    string `json:"user_id"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
		Author    string `json:"author"`
		Profile   string `json:"profile"`
	}

	err = row.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.Author, &comment.Profile)
	if err != nil {
		if err == sql.ErrNoRows {
			helper.RespondWithError(w, http.StatusNotFound, "Comment not found")
			return
		}
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, comment)
}
