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

	"social-network/internal/repository/midlweare"

	"github.com/google/uuid"
)

func Createpost(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
	userID, err := midlweare.AuthenticateUser(r)
	if err != nil {

		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	allowedUsers:= ""

	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	visability := strings.TrimSpace(r.FormValue("visibility"))
	 if visability == "private" {
		 allowedUsers := r.FormValue("allowed_users")
		 if allowedUsers == "" {
			 helper.RespondWithError(w, http.StatusBadRequest, "Allowed users must be provided for private posts")
			 return
		 }
	 }

	// if len(title) < 5 || len(title) > 50 || len(content) < 5 || len(content) > 500 {
	// 	helper.RespondWithError(w, http.StatusBadRequest, "Title and content must be between 5 and 50 characters")
	// 	return
	// }
   if visability == "private" {
		 allowedUsers = r.FormValue("allowed_users")
		 if allowedUsers == "" {
			 helper.RespondWithError(w, http.StatusBadRequest, "Allowed users must be provided for private posts")
			 return
		 }
	 }


	// if len(title) < 
	var imagePath string

	imageFile, _, err := r.FormFile("image")
	if err == nil {
		defer imageFile.Close()

		uploadDir := "../frontend/my-app/public/uploads"
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {

			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create upload directory")
			return
		}

		imagePath = fmt.Sprintf("uploads/%s.jpg", uuid.New().String()) // Keep the path relative for database storage

		out, err := os.Create(filepath.Join("../frontend/my-app/public", imagePath))
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to save image")
			return
		}
		defer out.Close()

		_, err = io.Copy(out, imageFile)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to save image")
			return
		}

	} else {
		imagePath = ""
	}

	postID := uuid.New().String()
	_, err = repository.Db.Exec(`
        INSERT INTO posts (id, user_id, title, content, visibility, image_path, canseperivite)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
		postID, userID, title, content, visability, imagePath, allowedUsers,
	)
	if err != nil {
		fmt.Println(err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "Post created successfully",
		"post_id": postID,
	})
}
