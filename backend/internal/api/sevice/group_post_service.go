package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"social-network/internal/repository"
	"social-network/internal/utils"

	"github.com/google/uuid"
)

func CreateGroupPostService(r *http.Request, userID string) (interface{}, error) {
	const maxFileSize = 1 * 1024 * 1024 * 1024 // 1GB

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		return nil, fmt.Errorf("group not found")
	}
	groupID := parts[3]

	// parse form (10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, fmt.Errorf("unable to parse form: %v", err)
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	// check membership
	isMember, err := repository.CheckUserInGroup(repository.Db, userID, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to check membership: %v", err)
	}
	if !isMember {
		return nil, fmt.Errorf("you are not a member of this group")
	}

	// handle file upload
	imagePath, err := handleGroupPostMedia(r, maxFileSize)
	if err != nil {
		return nil, err
	}

	postID := uuid.New().String()
	createdAt := time.Now().UTC()

	err = repository.InsertGroupPost(postID, userID, groupID, title, description, imagePath, createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %v", err)
	}

	newPost, err := repository.GetGroupPostByID(postID, userID)
	if err != nil {
		return nil, fmt.Errorf("post created but failed to fetch it: %v", err)
	}

	return newPost, nil
}

func handleGroupPostMedia(r *http.Request, maxSize int64) (string, error) {
	file, header, err := r.FormFile("image")
	if err != nil {
		return "", nil // no file uploaded
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".mp4": true, ".mov": true, ".avi": true}
	if !allowed[ext] {
		return "", fmt.Errorf("unsupported media format")
	}
	if header.Size > maxSize {
		return "", fmt.Errorf("file too large, max 1GB")
	}

	uploadDir := "../frontend/my-app/public/uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory")
	}

	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	imagePath := fmt.Sprintf("uploads/%s", fileName)
	out, err := os.Create(filepath.Join("../frontend/my-app/public", imagePath))
	if err != nil {
		return "", fmt.Errorf("failed to save file")
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("failed to write file")
	}
	return imagePath, nil
}

func GetAllGroupPostsService(r *http.Request, userID string) ([]utils.GroupPost, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		return nil, fmt.Errorf("group not found")
	}
	groupID := parts[3]

	isMember, err := repository.CheckUserInGroup(repository.Db,userID, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to check group membership")
	}
	if !isMember {
		return nil, fmt.Errorf("you are not a member of this group")
	}

	posts, err := repository.GetAllGroupPosts(groupID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group posts: %v", err)
	}

	return posts, nil
}
