package handlers

import (
	"net/http"

	service "social-network/app/api/service"
	"social-network/app/helper"
	"social-network/app/repository"
)

func GetStories(w http.ResponseWriter, r *http.Request) {
	authUserID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	stories, err := service.FetchStories(authUserID, repository.Db)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch stories")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, stories)
}
