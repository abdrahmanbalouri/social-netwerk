package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	service "social-network/internal/api/sevice"
	"social-network/internal/helper"
)

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Extract group ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusNotFound, "Group not found")
		return
	}
	groupID := parts[3]

	// Decode request body
	var payload struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DateTime    string `json:"dateTime"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Call service to create event
	eventID, err := service.CreateGroupEvent(userID, groupID, payload.Title, payload.Description, payload.DateTime)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Respond with created event info
	helper.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"id":          eventID,
		"title":       payload.Title,
		"description": payload.Description,
		"time":        payload.DateTime,
	})
}
