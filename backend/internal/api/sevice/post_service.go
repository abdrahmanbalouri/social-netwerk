package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"social-network/internal/repository"

	"github.com/google/uuid"
)

 //maxFileSize = 1 * 1024 * 1024 * 1024 // 1GB

func CreatePost(userID, title, content, visibility, allowedUsers string, fileHeader io.ReadCloser, filename string) (string, error) {
	// Validate title
	if len(title) < 2 || len(title) > 20 {
		return "", errors.New("title must be between 2 and 20 characters")
	}

	var allowed []string
	if visibility == "private" {
		allowed = strings.Split(allowedUsers, ",")
		if len(allowed) < 1 {
			return "", errors.New("allowed users must be provided for private posts")
		}
	}

	// Handle image upload
	var imagePath string
	if fileHeader != nil && filename != "" {
		defer fileHeader.Close()
		ext := strings.ToLower(filepath.Ext(filename))
		allowedExts := map[string]bool{
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
			".mp4": true, ".mov": true, ".avi": true,
		}
		if !allowedExts[ext] {
			return "", errors.New("unsupported media format")
		}
		uploadDir := "../frontend/my-app/public/uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create upload directory: %v", err)
		}

		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		imagePath = fmt.Sprintf("uploads/%s", filename)
		out, err := os.Create(filepath.Join("../frontend/my-app/public", imagePath))
		if err != nil {
			return "", fmt.Errorf("failed to save image: %v", err)
		}
		defer out.Close()
		if _, err := io.Copy(out, fileHeader); err != nil {
			return "", fmt.Errorf("failed to save image: %v", err)
		}
	}

	postID := uuid.New().String()
	post := repository.Post{
		ID:           postID,
		UserID:       userID,
		Title:        title,
		Content:      content,
		ImagePath:    imagePath,
		Visibility:   visibility,
		AllowedUsers: allowed,
	}

	if err := repository.InsertPost(repository.Db, post); err != nil {
		return "", err
	}

	if visibility == "private" {
		if err := repository.InsertAllowedUsers(repository.Db, postID, userID, allowed); err != nil {
			return "", err
		}
	}

	return postID, nil
}
