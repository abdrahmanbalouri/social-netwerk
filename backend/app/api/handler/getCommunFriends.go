package handlers

import (
	"net/http"

	"social-network/app/api/service"
	"social-network/app/helper"
)

func GetCommunFriends(w http.ResponseWriter, r *http.Request) {
	//
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, " method not allowed ")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	friends, statusCode, err := service.GetCommunFriends(userID)
	if err != nil {
		helper.RespondWithError(w, statusCode, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, friends)
}
