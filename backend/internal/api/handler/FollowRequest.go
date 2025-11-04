package handlers

import (
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/api/service"
)

func FollowRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method !=  http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, " method not allowed ")
		return
	}
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	result, err := service.GetFollowRequests(userID)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, result)
}
