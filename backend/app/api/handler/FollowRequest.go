package handlers

import (
	"net/http"

	"social-network/app/api/service"
	"social-network/app/helper"
)

func FollowRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, " method not allowed ")
		return
	}
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	result, statusCode, err := service.GetFollowRequests(userID)
if err != nil {
	helper.RespondWithError(w, statusCode, err.Error())
	return
}

	helper.RespondWithJSON(w, http.StatusOK, result)
}
