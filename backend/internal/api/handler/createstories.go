package handlers

import (
	"encoding/json"
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

type StoryRequest struct {
	Content string `json:"content"`
	BgColor string `json:"bg_color"`
}

func CreateStories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Limit upload size (5 MB max)
	const maxUploadSize = 5 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "File too large or invalid form")
		return
	}

	content := r.FormValue("content")
	imageFile, imageHeader, imgErr := r.FormFile("image")
	bgColor := r.FormValue("bg_color")
	if bgColor == "" {
		bgColor = "#000000"
	}
	content = strings.TrimSpace(content)

	if content == "" && imgErr != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Either content or image must be provided")
		return
	}
	if len(content) > 20 {
		helper.RespondWithError(w, http.StatusBadRequest, "content to large the top is 20")
		return

	}

	var imagePath string

	if imgErr == nil {
		defer imageFile.Close()

		// Validate MIME
		mimeType := imageHeader.Header.Get("Content-Type")
		allowedMimes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/webp": true,
			"image/gif":  true,
		}
		if !allowedMimes[mimeType] {
			helper.RespondWithError(w, http.StatusBadRequest, "Invalid image type")
			return
		}

		uploadDir := "../frontend/my-app/public/uploads/stories"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create upload directory")
			return
		}

		ext := filepath.Ext(imageHeader.Filename)
		if ext == "" {
			ext = ".jpg" // fallback
		}
		filename := uuid.New().String() + ext
		imagePath = "/uploads/stories/" + filename

		// Full server path
		serverPath := filepath.Join(uploadDir, filename)

		out, err := os.Create(serverPath)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to save image")
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, imageFile); err != nil {
			os.Remove(serverPath)
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to save image")
			return
		}
	} else {
		imagePath = ""
	}

	_, err = repository.Db.Exec(`
        INSERT INTO stories (user_id, content, image_url, bg_color)
        VALUES (?, ?, ?, ?)`,
		userID, content, imagePath, bgColor,
	)
	if err != nil {
		// If insert fails and image was saved, optionally delete it
		if imagePath != "" {
			os.Remove(filepath.Join("../frontend/my-app/public/uploads/stories", filepath.Base(imagePath)))
		}
		http.Error(w, "Failed to create story", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"image_url":  imagePath,
		"created_at": time.Now().Format("2006-01-02 15:04:05"),
	})
}
