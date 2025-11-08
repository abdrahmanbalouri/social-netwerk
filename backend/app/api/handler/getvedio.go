package handlers

import (
	"net/http"

	service "social-network/app/api/service"
	"social-network/app/helper"
	"social-network/app/repository"
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
