package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"social-network/internal/repository"

	"github.com/google/uuid"
)

const maxFileSize = 1 * 1024 * 1024 * 1024 // 1 GB

func CreateStory(userID, content, bgColor string, imageFile io.ReadCloser, filename string) (string, error) {
	content = strings.TrimSpace(content)
	if bgColor == "" {
		bgColor = "#000000"
	}

	if content == "" && imageFile == nil {
		return "", fmt.Errorf("either content or image must be provided")
	}

	if len(content) > 20 {
		return "", fmt.Errorf("content too long, max 20 characters")
	}

	var imagePath string
	if imageFile != nil && filename != "" {
		defer imageFile.Close()

		ext := filepath.Ext(filename)
		if ext == "" {
			ext = ".jpg"
		}
		imagePath = "/uploads/stories/" + uuid.New().String() + ext
		uploadDir := "../frontend/my-app/public/uploads/stories"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create upload directory: %v", err)
		}

		serverPath := filepath.Join(uploadDir, filepath.Base(imagePath))
		out, err := os.Create(serverPath)
		if err != nil {
			return "", fmt.Errorf("failed to save image: %v", err)
		}
		defer out.Close()

		if _, err := io.Copy(out, imageFile); err != nil {
			os.Remove(serverPath)
			return "", fmt.Errorf("failed to save image: %v", err)
		}
	}

	story := repository.Story{
		UserID:   userID,
		Content:  content,
		ImageURL: imagePath,
		BgColor:  bgColor,
	}

	if err := repository.InsertStory(repository.Db, story); err != nil {
		// Delete saved image if DB insert fails
		if imagePath != "" {
			os.Remove(filepath.Join("../frontend/my-app/public/uploads/stories", filepath.Base(imagePath)))
		}
		return "", err
	}

	return imagePath, nil
}
