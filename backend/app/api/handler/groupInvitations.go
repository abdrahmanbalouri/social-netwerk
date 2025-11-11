package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	service "social-network/app/api/service"
	"social-network/app/repository/model"
	"social-network/app/utils"

	"social-network/app/helper"
)

func GroupInvitationResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req utils.GroupResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "No valid session found")
		return
	}

	if err := service.ProcessGroupInvitationResponse(userID, req); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Invitation successfully processed",
	})
}

func GroupInvitationRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusNotFound, "Group not found")
		return
	}
	groupID := parts[3]

	var newInvitation utils.GroupInvitation
	if err := json.NewDecoder(r.Body).Decode(&newInvitation); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	response, err, errorStatus := service.HandleGroupInvitation(groupID, userID, newInvitation)
	if err != nil {
		helper.RespondWithError(w, errorStatus, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, response)
}

func FetchJoinRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusNotFound, "Group not found")
		return
	}
	groupID := parts[3]

	joinRequests, err := model.FetchJoinRequests(userID, groupID)
	if err != nil {
		http.Error(w, "Error fetching join requests: "+err.Error(), http.StatusInternalServerError)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, joinRequests)
}

func FetchFriendsForGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	groupID := parts[3]

	users, err := model.FetchFriendsNotInGroup(userID, groupID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, users)
}
