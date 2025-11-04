package handlers

import (
	"net/http"
	"social-network/internal/api/service"
	"social-network/internal/helper"
)

func FollowHandler(w http.ResponseWriter, r *http.Request) {


if r.Method !=  http.MethodPost {
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

	response, err := service.ToggleFollow(userID, targetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, response)
}
