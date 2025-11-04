package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	service "social-network/internal/api/sevice"
	"social-network/internal/helper"
	"social-network/internal/repository"
)

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	authUserID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	postID := parts[3]

	postResp, err := service.FetchPost(repository.Db, postID, authUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			helper.RespondWithError(w, http.StatusNotFound, "Post not found")
			return
		}
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch post")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, postResp)
}
