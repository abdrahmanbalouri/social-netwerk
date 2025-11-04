package handlers

import (
	"net/http"

	service "social-network/internal/api/sevice"
	"social-network/internal/helper"
	"social-network/internal/repository"
)

func GetVideoHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	posts, err := service.FetchVideoPosts(userID, repository.Db)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch videos")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, posts)
}
