package handlers

import (
	"net/http"

	service "social-network/app/api/service"
	"social-network/app/helper"
)

func Getfollowers(w http.ResponseWriter, r *http.Request) {
	// Assume token is sent in Authorization header
	UserId, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	followers, err := service.GetFollowersService(UserId)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, followers)
}
