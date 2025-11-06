package handlers

import (
	"fmt"
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository/model"
)

func GroupeInvitation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WAST HAD LBATAAAAAL")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	fmt.Println("UserID :::::", userID)

	groupInvitations, err := model.FetchGroupInvitations(userID)
	if err != nil {
		http.Error(w, "Error fetching group invitations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, groupInvitations)
}
