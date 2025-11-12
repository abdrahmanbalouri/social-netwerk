package handlers

import (
	"net/http"
	"time"

	service "social-network/app/api/service"
	"social-network/app/helper"
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

// {
//     "grpId": "b5212293-b4db-40d6-b0e0-7f68a143d2b8"
// }
