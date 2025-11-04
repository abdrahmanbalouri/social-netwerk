package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"social-network/internal/helper"
	"social-network/internal/repository"
)

//const maxFileSize = 1 * 1024 * 1024 * 1024 // 1 GB

func CreateComment(userID, postID, content, whatis, groupID string, mediaFileHeader map[string]interface{}) (string, string, error) {
	if postID == "" {
		return "", "", errors.New("missing post ID")
	}
	if content == "" && mediaFileHeader == nil {
		return "", "", errors.New("either content or media is required")
	}
	if len(content) > 300 {
		return "", "", errors.New("content too long")
	}

	sanitizedContent := helper.Skip(content)

	if whatis == "groups" {
		err := helper.CheckUserInGroup(userID, groupID)
		if err != nil {
			return "", "", errors.New("user not in group")
		}
	}

	// Handle media upload
	var mediaPath string
	if mediaFileHeader != nil {
		file := mediaFileHeader["file"].(io.Reader)
		filename := mediaFileHeader["filename"].(string)

		ext := strings.ToLower(filepath.Ext(filename))
		allowedExts := map[string]bool{
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
			".mp4": true, ".mov": true, ".avi": true,
		}
		if !allowedExts[ext] {
			return "", "", errors.New("unsupported media format")
		}

		if mediaFileHeader["size"].(int64) > maxFileSize {
			return "", "", errors.New("file too large, max 1 GB")
		}

		uploadDir := "../frontend/my-app/public/uploads"
		os.MkdirAll(uploadDir, os.ModePerm)

		filename = fmt.Sprintf("%s%s", uuid.New().String(), ext)
		mediaPath = fmt.Sprintf("uploads/%s", filename)
		out, err := os.Create(filepath.Join("../frontend/my-app/public", mediaPath))
		if err != nil {
			return "", "", errors.New("failed to save media file")
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			return "", "", errors.New("failed to save media file")
		}
	}

	commentID := uuid.New().String()
	comment := repository.Comment{
		ID:        commentID,
		PostID:    postID,
		UserID:    userID,
		Content:   sanitizedContent,
		MediaPath: mediaPath,
		GroupID:   groupID,
		IsGroup:   whatis == "groups",
	}

	// Save to database
	err := repository.InsertComment(repository.Db, comment)
	if err != nil {
		return "", "", err
	}

	return commentID, mediaPath, nil
}
