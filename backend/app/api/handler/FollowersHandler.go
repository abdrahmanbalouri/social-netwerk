package handlers

import (
	"net/http"

	"social-network/app/api/service"
	"social-network/app/helper"
)

func FollowersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, " method not allowed ")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	targetID := r.URL.Query().Get("id")
	if targetID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "missing target user ID")
		return
	}

	followers, err := service.GetFollowers(userID, targetID)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, followers)
}
