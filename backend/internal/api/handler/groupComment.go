package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"social-network/internal/helper"
	"social-network/internal/repository"

	"github.com/google/uuid"
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
	if r.Method != "POST" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = r.ParseMultipartForm(20 << 20)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	postID := strings.TrimSpace(r.FormValue("postId"))
	content := strings.TrimSpace(r.FormValue("content"))
	if postID == "" || content == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}
	if len(content) < 3 || len(content) > 300 {
		helper.RespondWithError(w, http.StatusBadRequest, "Content length invalid")
		return
	}

	sanitizedContent := helper.Skip(content)

	row := repository.Db.QueryRow(`SELECT content FROM group_posts WHERE id = ?`, postID)
	var exists string
	err = row.Scan(&exists)
	if err == sql.ErrNoRows {
		helper.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	} else if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	var mediaPath string
	file, header, err := r.FormFile("media")
	if err == nil {
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		allowedExts := map[string]bool{
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
			".mp4": true, ".mov": true, ".avi": true,
		}
		if !allowedExts[ext] {
			helper.RespondWithError(w, http.StatusBadRequest, "Unsupported media format")
			return
		}

		uploadDir := "../frontend/my-app/public/uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create upload directory")
			return
		}

		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		mediaPath = fmt.Sprintf("uploads/%s", filename)

		out, err := os.Create(filepath.Join("../frontend/my-app/public", mediaPath))
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to save media file")
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to save media file")
			return
		}
	} else {
		mediaPath = ""
	}

	commentID := uuid.New().String()
	fmt.Println(commentID, "------------")

	_, err = repository.Db.Exec(`
		INSERT INTO comments (id, post_id, user_id, content, media_path)
		VALUES (?, ?, ?, ?, ?)`,
		commentID, postID, userID, sanitizedContent, mediaPath)
	if err != nil {
		fmt.Println(err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"message":    "Comment created successfully",
		"comment_id": commentID,
		"media":      mediaPath,
	})
}

func GetCommentGroup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("INSIDE GET COMMENTS")
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	
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
