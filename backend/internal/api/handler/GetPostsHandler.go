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

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
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

	var post struct {
		ID            string    `json:"id"`
		Title         string    `json:"title"`
		Content       string    `json:"content"`
		Image_path    string    `json:"image_path"`
		CreatedAt     time.Time `json:"created_at"`
		Author        string    `json:"author"`
		Profile       string    `json:"profile"`
		CommentsCount int       `json:"comments_count"`
	}

	err := repository.Db.QueryRow(`
		SELECT 
			p.id, 
			p.title, 
			p.content, 
			p.image_path, 
			p.created_at, 
			u.nickname,
			u.image AS profile,
			COUNT(c.id) AS comments_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE p.id = ?
		GROUP BY p.id, p.title, p.content, p.image_path, p.created_at, u.nickname, u.image`, postID).Scan(
		&post.ID, &post.Title, &post.Content, &post.Image_path, &post.CreatedAt, &post.Author, &post.Profile, &post.CommentsCount)
	if err != nil {
		if err == sql.ErrNoRows {
			helper.RespondWithError(w, http.StatusNotFound, "Post not found")
			return
		}
		fmt.Printf("Error querying post: %v\n", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch post")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, post)
}