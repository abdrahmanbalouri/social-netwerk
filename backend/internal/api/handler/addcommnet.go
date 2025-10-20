package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"social-network/internal/helper"
	"social-network/internal/repository"

	"github.com/google/uuid"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	err = r.ParseMultipartForm(20 << 20) // 20 MB max
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

	row := repository.Db.QueryRow(`SELECT content FROM posts WHERE id = ?`, postID)
	var exists string
	err = row.Scan(&exists)
	if err == sql.ErrNoRows {
		helper.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	} else if err != nil {
		fmt.Println("1")
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
			fmt.Println("22")
			return
		}

		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		mediaPath = fmt.Sprintf("uploads/%s", filename)

		out, err := os.Create(filepath.Join("../frontend/my-app/public", mediaPath))
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to save media file")
			fmt.Println("33")
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to save media file")
			fmt.Println("44")
			return
		}
	} else {
		mediaPath = ""
	}

	commentID := uuid.New().String()

	_, err = repository.Db.Exec(`
		INSERT INTO comments (id, post_id, user_id, content, media_path)
		VALUES (?, ?, ?, ?, ?)`,
		commentID, postID, userID, sanitizedContent, mediaPath)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create comment")
		fmt.Println("55")
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"message":    "Comment created successfully",
		"comment_id": commentID,
		"media":      mediaPath,
	})
}
