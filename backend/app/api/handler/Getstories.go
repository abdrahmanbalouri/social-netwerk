package handlers

import (
	"net/http"

	service "social-network/app/api/service"
	"social-network/app/helper"
	"social-network/pkg/db/sqlite"
)

func GetStories(w http.ResponseWriter, r *http.Request) {

	if r.Method !=  http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, " method not allowed ")
		return
}

	authUserID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	stories, err := service.FetchStories(authUserID, sqlite.Db)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch stories")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, stories)
}
