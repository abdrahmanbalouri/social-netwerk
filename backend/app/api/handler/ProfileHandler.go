package handlers

import (
	"net/http"

	"social-network/app/helper"
	"social-network/app/repository/profile"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, " method not allowed ")
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
	profileData, err := profile.GetProfileData(targetUserID, currentUserID, w)
	if err != nil {
		http.Error(w, "Failed to get profile data", http.StatusInternalServerError)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, profileData)
}
