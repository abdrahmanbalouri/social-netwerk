package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"social-network/internal/helper"
	"social-network/internal/repository"
	"social-network/internal/repository/middleware"

	"github.com/google/uuid"
)

func Createpost(w http.ResponseWriter, r *http.Request) {
const maxFileSize = 1 * 1024 * 1024 * 1024 // 1 gb bazaf yak

	if r.Method != "POST" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	userID, err := middleware.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	title := helper.Skip(strings.TrimSpace(r.FormValue("title")))
	content := helper.Skip(strings.TrimSpace(r.FormValue("content")))
	visibility := strings.TrimSpace(r.FormValue("visibility"))
	allowedUsers := strings.TrimSpace(r.FormValue("allowed_users"))
	if  len(title) > 20 {
		helper.RespondWithError(w, http.StatusBadRequest, "title bitwen 2 and  20")
		return 
	}
	users := len(strings.Split(allowedUsers, ","))
	if visibility == "private" && users == 1 {
		helper.RespondWithError(w, http.StatusBadRequest, "Allowed users must be provided for private posts")
		return
	}

	// Handle image upload
	var imagePath string
	imageFile, header, err := r.FormFile("image")
	if err == nil {
		defer imageFile.Close()
		ext := strings.ToLower(filepath.Ext(header.Filename))
		allowedExts := map[string]bool{
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
			".mp4": true, ".mov": true, ".avi": true,
		}
		if !allowedExts[ext] {
			helper.RespondWithError(w, http.StatusBadRequest, "Unsupported media format")
			return
		}
		 if header.Size > maxFileSize {
         helper.RespondWithError(w, http.StatusBadRequest, "File too large, max 1 GB")
        return
    }
		uploadDir := "../frontend/my-app/public/uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create upload directory")
			return
		}
		filname := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		imagePath = fmt.Sprintf("uploads/%s", filname)

		out, err := os.Create(filepath.Join("../frontend/my-app/public", imagePath))
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to save image")
			return
		}
		defer out.Close()
		if _, err := io.Copy(out, imageFile); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to save image")
			return
		}
	} else {
		imagePath = ""
	}
	postID := uuid.New().String()

	_, err = repository.Db.Exec(`
		INSERT INTO posts (id, user_id, title, content, image_path, visibility, canseperivite)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		postID, userID, title, content, imagePath, visibility, allowedUsers,
	)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	// Insert allowed users for private posts
	if visibility == "private" {
		userIDs := strings.Split(allowedUsers, ",")
		for _, uid := range userIDs {
			uid = strings.TrimSpace(uid)
			uid = strings.ReplaceAll(uid, `"`, "") // remove any quotes

			if uid == "" {
				continue
			}
			// Check if user exists to avoid FK errors
			var exists int
			err := repository.Db.QueryRow(`SELECT 1 FROM users WHERE id = ?`, uid).Scan(&exists)
			if err != nil {
				continue
			}

			_, err = repository.Db.Exec(`
				INSERT INTO allowed_followers (user_id, post_id, allowed_user_id)
				VALUES (?, ?, ?)
			`, userID, postID, uid)
			if err != nil {
				continue
			}
		}
	}

	helper.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "Post created successfully",
		"post_id": postID,
	})
}
