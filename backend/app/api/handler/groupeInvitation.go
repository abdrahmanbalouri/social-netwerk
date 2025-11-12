package handlers

import (
	"net/http"

	"social-network/app/helper"
	"social-network/app/repository/model"
)

func GroupeInvitation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	groupInvitations, err := model.FetchGroupInvitations(userID)
	if err != nil {
		http.Error(w, "Error fetching group invitations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, groupInvitations)
}
