package handlers

import (
	"net/http"

	"social-network/app/api/service"
	"social-network/app/helper"
)

func FollowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusUnauthorized, "unauthorised")
		return
	}
	w.Header().Set("Content-Type", "application/json")

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

	response, statusCode, err := service.ToggleFollow(userID, targetID)
	if err != nil {
		helper.RespondWithError(w, statusCode, err.Error())

		return
	}

	helper.RespondWithJSON(w, http.StatusOK, response)
}
