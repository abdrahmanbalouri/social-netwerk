package handlers

import (
	"encoding/json"
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

type GroupPost struct {
	GrpID   string `json:"grpId"`
	Title   string `json:"title"`
	Content string `json:"content"`
	ImgPath string `json:"imgPath"`
}

func CreatePostGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	// check for inputs
	var newPost GroupPost
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}
	if len(strings.TrimSpace(newPost.Title)) == 0 || len(strings.TrimSpace(newPost.Content)) == 0 {
		helper.RespondWithError(w, http.StatusBadRequest, "Title and description are required")
		return
	}

	// start transaction
	tx, TransErr := repository.Db.Begin()
	if TransErr != nil {
		fmt.Println("Error starting the transaction")
		helper.RespondWithError(w, http.StatusInternalServerError, "Error starting the transaction")
		return
	}

	defer tx.Rollback()

	// checl auth
	userID, IdErr := helper.AuthenticateUser(r)
	if IdErr != nil {
		fmt.Println("Error getting the user's id")
		helper.RespondWithError(w, http.StatusInternalServerError, IdErr.Error())
		return
	}

	// check for membership of the user
	var isMember bool
	query := `SELECT EXISTS (SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?)`
	err := tx.QueryRow(query, userID, newPost.GrpID).Scan(&isMember)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to check group membership")
		return
	}

	if !isMember {
		fmt.Println("The user is not a member of the group")
		helper.RespondWithError(w, http.StatusUnauthorized, "You are not a member of this group")
		return
	}
	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))

	var imagePath string
	imageFile, _, err := r.FormFile("image")
	if err == nil {
		defer imageFile.Close()

		uploadDir := "../frontend/my-app/public/uploads"
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			fmt.Println("Failed to create upload directory")
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create upload directory")
			return
		}

		imagePath = fmt.Sprintf("uploads/%s.jpg", uuid.New().String()) // Keep the path relative for database storage
		out, err := os.Create(filepath.Join("../frontend/my-app/public", imagePath))
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to save image")
			fmt.Println("Failed to save image :", err)
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

	postID := helper.GenerateUUID()
	createdAt := time.Now().UTC()
	_, err = repository.Db.Exec(`
        INSERT INTO posts (id, user_id, group_id, title, content, image_path, created_at)
        VALUES (?, ?, ?, ?, ?)`,
		postID, userID, title, content, imagePath,
	)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	if err := tx.Commit(); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	response := map[string]string{
		"message": "post created successfully for group",
	}
	helper.RespondWithJSON(w, http.StatusOK, response)
}
