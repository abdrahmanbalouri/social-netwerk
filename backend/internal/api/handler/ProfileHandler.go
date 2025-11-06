package handlers

import (
	"net/http"

	"social-network/internal/helper"
	profilego "social-network/internal/repository/profile"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	targetUserID := r.URL.Query().Get("userId")
	currentUserID, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if targetUserID == "0" || targetUserID == "" {
		targetUserID = currentUserID
	}
	profileData, err := profilego.GetProfileData(targetUserID, currentUserID, w)
	if err != nil {
		http.Error(w, "Failed to get profile data", http.StatusInternalServerError)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, profileData)
}
