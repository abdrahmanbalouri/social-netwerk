package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"social-network/internal/repository"
	"social-network/internal/repository/model"
)

func UpdateUserProfile(userID, displayName, privacy string, coverFile io.ReadCloser, coverName string, avatarFile io.ReadCloser, avatarName string) (string, string, error) {
	// Fetch current files
	oldCover, oldAvatar, err := model.GetProfileFiles(repository.Db, userID)
	if err != nil {
		return "", "", err
	}

	uploadDir := "../frontend/my-app/public/uploads"

	// Handle cover file
	coverFilename := oldCover
	if coverFile != nil && coverName != "" {
		defer coverFile.Close()
		coverPath := filepath.Join(uploadDir, coverName)
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return "", "", fmt.Errorf("failed to create directory: %v", err)
		}
		out, err := os.Create(coverPath)
		if err != nil {
			return "", "", err
		}
		defer out.Close()
		if _, err := io.Copy(out, coverFile); err != nil {
			return "", "", err
		}
		coverFilename = coverName
	}

	// Handle avatar file
	avatarFilename := oldAvatar
	if avatarFile != nil && avatarName != "" {
		defer avatarFile.Close()
		avatarPath := filepath.Join(uploadDir, avatarName)
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return "", "", fmt.Errorf("failed to create directory: %v", err)
		}
		out, err := os.Create(avatarPath)
		if err != nil {
			return "", "", err
		}
		defer out.Close()
		if _, err := io.Copy(out, avatarFile); err != nil {
			return "", "", err
		}
		avatarFilename = avatarName
	}

	profile := model.Profile{
		UserID:      userID,
		DisplayName: displayName,
		Privacy:     privacy,
		Cover:       coverFilename,
		Avatar:      avatarFilename,
	}

	if err := model.UpdateProfile(repository.Db, profile); err != nil {
		return "", "", err
	}

	return coverFilename, avatarFilename, nil
}
