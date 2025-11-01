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

type PostData struct {
	GrpID   string `json:"grpId"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type FetchPost struct {
	GrpID string `json:"grpId"`
}

type Post struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	ImagePath    string    `json:"image_path"`
	CreatedAt    time.Time `json:"created_at"`
	Author       string    `json:"author"`
	Profile      string    `json:"profile"`
	Like         int       `json:"likeCount"`
	LikedByUSer  int       `json:"liked_by_user"`
	CommentCount int       `json:"comment_count"`
}

func CreatePostGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// get user id
	userID, IdErr := helper.AuthenticateUser(r)
	if IdErr != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, IdErr.Error())
		return
	}

	// parse Data
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Unable to parse form")
		return
	}
	data := r.FormValue("postData")
	if data == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing JSON form field")
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusNotFound, "Group not found")
		return
	}
	GrpID := parts[3]
	var postData PostData
	err = json.Unmarshal([]byte(data), &postData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// check for membership of the user
	var isMember bool
	query := `SELECT EXISTS (SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?)`
	err = repository.Db.QueryRow(query, userID, GrpID).Scan(&isMember)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to check group membership")
		return
	}
	if !isMember {
		helper.RespondWithError(w, http.StatusUnauthorized, "You are not a member of this group")
		return
	}

	// image part
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

	// insert new post into posts tbale
	postID := helper.GenerateUUID()
	createdAt := time.Now().UTC()
	_, err = repository.Db.Exec(`
        INSERT INTO group_posts (id, user_id, group_id, title, content, image_path, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
		postID, userID, GrpID, postData.Title, postData.Content, imagePath, createdAt,
	)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	var newPost struct {
		ID            string    `json:"id"`
		UserID        string    `json:"user_id"`
		Title         string    `json:"title"`
		Content       string    `json:"content"`
		ImagePath     string    `json:"image_path"`
		CreatedAt     time.Time `json:"created_at"`
		Nickname      string    `json:"nickname"`
		Profile       string    `json:"profile"`
		LikeCount     int       `json:"like_count"`
		LikedByUser   int       `json:"liked_by_user"`
		CommentsCount int       `json:"comments_count"`
	}

	queryNewPost := `
	SELECT  
	gp.id, gp.user_id, gp.title, gp.content, gp.image_path, gp.created_at, 
	u.nickname, u.image AS profile, 
	COUNT(DISTINCT l.id) AS like_count, 
	COUNT(DISTINCT CASE WHEN l.user_id = ? THEN l.id END) AS liked_by_user, 
	COUNT(DISTINCT c.id) AS comments_count 
	FROM group_posts gp 
	JOIN users u ON gp.user_id = u.id 
	LEFT JOIN likes l ON gp.id = l.liked_item_id AND l.liked_item_type = 'post' 
	LEFT JOIN comments c ON gp.id = c.post_id 
	GROUP BY gp.id, gp.user_id, gp.title, gp.content, gp.image_path, gp.created_at, u.nickname, u.image
	ORDER BY gp.created_at DESC
	LIMIT 1;
	`

	err = repository.Db.QueryRow(queryNewPost, userID, postID).Scan(
		&newPost.ID, &newPost.UserID, &newPost.Title, &newPost.Content, &newPost.ImagePath,
		&newPost.CreatedAt, &newPost.Nickname, &newPost.Profile,
		&newPost.LikeCount, &newPost.LikedByUser, &newPost.CommentsCount,
	)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Post created but failed to fetch it")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, newPost)
}

func GetAllPostsGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusNotFound, "Group not found")
		return
	}
	GrpId := parts[3]

	// Get user ID
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication failed")
		return
	}

	// Check for user's membership
	var isMember bool
	query := `SELECT EXISTS (SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?)`
	if err := repository.Db.QueryRow(query, userID, GrpId).Scan(&isMember); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to check group membership")
		return
	}
	if !isMember {
		helper.RespondWithError(w, http.StatusUnauthorized, "You are not a member of this group")
		return
	}

	// Fetch all the posts of this group
	// query = `SELECT id, title, content, image_path, created_at, user_id FROM posts WHERE group_id = ?`
	query = `
	SELECT 
    gp.id, 
    gp.user_id, 
    gp.title, 
    gp.content, 
    gp.image_path, 
    gp.created_at, 
  u.first_name,u.last_name,
    u.image AS profile,
    COUNT(DISTINCT l.id) AS like_count,
    COUNT(DISTINCT CASE WHEN l.user_id = ? THEN l.id END) AS liked_by_user,
    COUNT(DISTINCT c.id) AS comments_count
	FROM group_posts gp
	JOIN users u ON gp.user_id = u.id
	LEFT JOIN likes l ON gp.id = l.liked_item_id AND l.liked_item_type = 'post'
	LEFT JOIN comments c ON gp.id = c.post_id
	WHERE gp.group_id = ?
	GROUP BY 
    gp.id, gp.user_id, gp.title, gp.content, gp.image_path, gp.created_at, 
     u.first_name,u.last_name, u.image
	ORDER BY gp.created_at DESC;
`
	rows, err := repository.Db.Query(query, userID, GrpId)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get posts")
		return
	}
	defer rows.Close()

	var postsJson []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.ImagePath, &p.CreatedAt, &p.Author, &p.Profile, &p.Like, &p.LikedByUSer, &p.CommentCount)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to scan posts")
			return
		}
		postsJson = append(postsJson, p)
	}

	// Return the posts as a JSON response
	helper.RespondWithJSON(w, http.StatusOK, postsJson)
}

func GetPostGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusNotFound, "Group not found")
		return
	}
	GrpId := parts[3]

	// Get user ID
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication failed")
		return
	}

	// Check for user's membership
	var isMember bool
	query := `SELECT EXISTS (SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?)`
	if err := repository.Db.QueryRow(query, userID, GrpId).Scan(&isMember); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to check group membership")
		return
	}
	if !isMember {
		helper.RespondWithError(w, http.StatusUnauthorized, "You are not a member of this group")
		return
	}

	// Fetch all the posts of this group
	// query = `SELECT id, title, content, image_path, created_at, user_id FROM posts WHERE group_id = ?`
	query = `
	SELECT 
    gp.id, 
    gp.user_id, 
    gp.title, 
    gp.content, 
    gp.image_path, 
    gp.created_at, 
    u.nickname,
    u.image AS profile,
    COUNT(DISTINCT l.id) AS like_count,
    COUNT(DISTINCT CASE WHEN l.user_id = ? THEN l.id END) AS liked_by_user,
    COUNT(DISTINCT c.id) AS comments_count
	FROM group_posts gp
	JOIN users u ON gp.user_id = u.id
	LEFT JOIN likes l ON gp.id = l.liked_item_id AND l.liked_item_type = 'post'
	LEFT JOIN comments c ON gp.id = c.post_id
	WHERE gp.group_id = ?
	GROUP BY 
    gp.id, gp.user_id, gp.title, gp.content, gp.image_path, gp.created_at, 
    u.nickname, u.image
	ORDER BY gp.created_at DESC
	LIMIT 1;
`
	rows, err := repository.Db.Query(query, userID, GrpId)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get posts")
		return
	}
	defer rows.Close()

	var p Post
	err = rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.ImagePath, &p.CreatedAt, &p.Author, &p.Profile, &p.Like, &p.LikedByUSer, &p.CommentCount)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to scan posts")
		return
	}

	// Return the posts as a JSON response
	helper.RespondWithJSON(w, http.StatusOK, p)
}

// {
//     "grpId": "b5212293-b4db-40d6-b0e0-7f68a143d2b8"
// }
