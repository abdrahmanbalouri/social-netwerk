package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"social-network/app/helper"
	"social-network/app/repository/model"
	"social-network/app/utils"
	"social-network/pkg/db/sqlite"
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

	title := strings.TrimSpace(r.FormValue("title"))
	description := strings.TrimSpace(r.FormValue("description"))

	if len(title) > 20 || len(description) > 40 {
		return nil, fmt.Errorf("title or description to bigg")
	}

	// check membership
	isMember, err := model.CheckUserInGroup(sqlite.Db, userID, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to check membership: %v", err)
	}
	if !isMember {
		return nil, fmt.Errorf("you are not a member of this group")
	}

	// handle file upload
	imagePath, err := handleGroupPostMedia(r, maxFileSize, title, description)
	if err != nil {
		return nil, err
	}

	postID := helper.GenerateUUID().String()
	createdAt := time.Now().UTC()

	err = model.InsertGroupPost(postID, userID, groupID, title, description, imagePath, createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %v", err)
	}

	newPost, err := model.GetGroupPostByID1(postID, userID)
	if err != nil {
		return nil, fmt.Errorf("post created but failed to fetch it: %v", err)
	}

	return newPost, nil
}

func handleGroupPostMedia(r *http.Request, maxSize int64, title string, descreption string) (string, error) {
	file, header, err := r.FormFile("image")
	if err != nil {
		return "", nil // no file uploaded
	}
	defer file.Close()
if title == "" && descreption == "" && file == nil {
    return "", fmt.Errorf("please provide at least one field")
}


	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".mp4": true, ".mov": true}
	if !allowed[ext] {
		return "", fmt.Errorf("unsupported media format")
	}
	if header.Size > maxSize {
		return "", fmt.Errorf("file too large, max 1GB")
	}

	uploadDir := "../frontend/public/uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory")
	}

	fileName := fmt.Sprintf("%s%s", helper.GenerateUUID().String(), ext)
	imagePath := fmt.Sprintf("uploads/%s", fileName)
	out, err := os.Create(filepath.Join("../frontend/public", imagePath))
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

	isMember, err := model.CheckUserInGroup(sqlite.Db, userID, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to check group membership")
	}
	if !isMember {
		return nil, fmt.Errorf("you are not a member of this group")
	}

	posts, err := model.GetAllGroupPosts(groupID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group posts: %v", err)
	}

	return posts, nil
}

func GetGroupPost(postID, userID, groupID string) (map[string]interface{}, error) {
	if err := helper.CheckUserInGroup(userID, groupID); err != nil {
		return nil, errors.New("user not in group")
	}

	post, err := model.GetGroupPostByID(sqlite.Db, userID, postID)
	if err != nil {
		return nil, err
	}

	postMap := map[string]interface{}{
		"id":             post.ID,
		"user_id":        post.UserID,
		"title":          post.Title,
		"content":        post.Content,
		"image_path":     post.ImagePath,
		"created_at":     post.CreatedAt,
		"first_name":     post.FirstName,
		"last_name":      post.LastName,
		"profile":        post.Profile,
		"like":           post.LikeCount,
		"liked_by_user":  post.LikedByUser,
		"comments_count": post.CommentsCount,
	}

	return postMap, nil
}
