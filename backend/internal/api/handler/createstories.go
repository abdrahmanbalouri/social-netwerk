package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"social-network/internal/helper"
	"social-network/internal/repository"

	"github.com/google/uuid" // <-- add this
)

type StoryRequest struct {
	Content string `json:"content"` // fallback if sent as JSON, but we use form
	BgColor string `json:"bg_color"`
}

// CreateStories handles multipart/form-data: content (text), bg_color (text), image (file optional)
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
	const maxUploadSize = 5 << 20 // 5 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "File too large or invalid form")
		return
	}

	// Get form fields
	content := r.FormValue("content")
	imageFile, imageHeader, imgErr := r.FormFile("image")
	bgColor := r.FormValue("bg_color")
	if bgColor == "" {
		bgColor = "#000000"
	}

	// Validate: at least content or image
	if content == "" && imgErr != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Either content or image must be provided")
		return
	}

	var imagePath string // public URL like "/uploads/stories/abc123.jpg"

	// Handle image if provided
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

		// Create upload dir (relative to backend; adjust if needed)
		uploadDir := "../frontend/my-app/public/uploads/stories"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create upload directory")
			return
		}

		// Generate unique filename (preserve original ext if possible)
		ext := filepath.Ext(imageHeader.Filename)
		if ext == "" {
			ext = ".jpg" // fallback
		}
		filename := uuid.New().String() + ext
		imagePath = "/uploads/stories/" + filename // public URL

		// Full server path
		serverPath := filepath.Join(uploadDir, filename)

		out, err := os.Create(serverPath)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to save image")
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, imageFile); err != nil {
			os.Remove(serverPath) // cleanup partial file
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to save image")
			return
		}
	} else {
		// No image provided (optional)
		imagePath = ""
	}

	// Insert into DB
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

	// Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"image_url":  imagePath,
		"created_at": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// GetStories remains the same as your provided code...
// (add the full GetStories func here if needed, with the SQL fix as in previous messages)
