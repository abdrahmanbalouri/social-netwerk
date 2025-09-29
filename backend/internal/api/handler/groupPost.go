package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ImagePath string    `json:"image_path"`
	CreatedAt time.Time `json:"created_at"`
	AuthorID  string    `json:"author_id"`
}

func CreatePostGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// get user id
	userID, IdErr := helper.AuthenticateUser(r)
	if IdErr != nil {
		fmt.Println("Error getting the user's id")
		helper.RespondWithError(w, http.StatusInternalServerError, IdErr.Error())
		return
	}

	// parse Data
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		fmt.Println("Unable to parse form")
		helper.RespondWithError(w, http.StatusBadRequest, "Unable to parse form")
		return
	}
	data := r.FormValue("postData")
	if data == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing JSON form field")
		return
	}
	var postData PostData
	err = json.Unmarshal([]byte(data), &postData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// check for membership of the user
	var isMember bool
	query := `SELECT EXISTS (SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?)`
	err = repository.Db.QueryRow(query, userID, postData.GrpID).Scan(&isMember)
	if err != nil {
		fmt.Println("Failed to check group membership")
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to check group membership")
		return
	}
	if !isMember {
		fmt.Println("The user is not a member of the group")
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
			fmt.Println("Failed to create upload directory")
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create upload directory")
			return
		}

		imagePath = fmt.Sprintf("uploads/%s.jpg", uuid.New().String()) // Keep the path relative for database storage
		out, err := os.Create(filepath.Join("../frontend/my-app/public", imagePath))
		if err != nil {
			fmt.Println("Failed to save image :", err)
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
        INSERT INTO posts (id, user_id, group_id, title, content, image_path, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
		postID, userID, postData.GrpID, postData.Title, postData.Content, imagePath, createdAt,
	)
	if err != nil {
		fmt.Println("Failed to create post (groups)")
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	response := map[string]string{
		"message": "post created successfully for group",
	}
	helper.RespondWithJSON(w, http.StatusOK, response)
}

func GetPostGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	var newPost FetchPost
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Get user ID
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication failed")
		return
	}

	// Check for user's membership
	var isMember bool
	query := `SELECT EXISTS (SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?)`
	if err := repository.Db.QueryRow(query, userID, newPost.GrpID).Scan(&isMember); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to check group membership")
		return
	}
	if !isMember {
		helper.RespondWithError(w, http.StatusUnauthorized, "You are not a member of this group")
		return
	}

	// Fetch all the posts of this group
	query = `SELECT id, title, content, image_path, created_at, user_id FROM posts WHERE group_id = ?`
	rows, err := repository.Db.Query(query, newPost.GrpID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get posts")
		return
	}
	defer rows.Close()

	var postsJson []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.ImagePath, &p.CreatedAt, &p.AuthorID);
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to scan posts")
			return
		}
		postsJson = append(postsJson, p)
	}

	// Return the posts as a JSON response
	helper.RespondWithJSON(w, http.StatusOK, postsJson)
}


// {
//     "grpId": "b5212293-b4db-40d6-b0e0-7f68a143d2b8"
// }