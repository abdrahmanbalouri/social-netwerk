package handlers

import (
	"net/http"
	"strings"
	"time"

	service "social-network/app/api/service"
	"social-network/app/helper"
	"social-network/app/repository"
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
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	ImagePath     string    `json:"image_path"`
	CreatedAt     time.Time `json:"created_at"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Profile       string    `json:"profile"`
	Like          int       `json:"like"`
	LikedByUSer   int       `json:"liked_by_user"`
	CommentsCount int       `json:"comments_count"`
}

func CreatePostGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	post, err := service.CreateGroupPostService(r, userID)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, post)
}

func GetAllPostsGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication failed")
		return
	}

	posts, err := service.GetAllGroupPostsService(r, userID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, posts)
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
	err = rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.ImagePath, &p.CreatedAt, &p.FirstName, &p.LastName, &p.Profile, &p.Like, &p.LikedByUSer, &p.CommentsCount)
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
