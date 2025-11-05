package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository/model"
)

func GetGroupMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// Authenticate user

	currentUserID, err := helper.AuthenticateUser(r)
	if err != nil {
		fmt.Println("Authentication error:", err)
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	// Parse query parameters
	groupId := r.URL.Query().Get("groupId")
	if groupId == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing groupId parameter")
		return
	}

	messages, err := model.GetGroupMessages(currentUserID, groupId)
	if err == sql.ErrNoRows {
		helper.RespondWithError(w, http.StatusNotFound, "No messages found")
		return
	} else if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Error fetching messages: "+err.Error())
		return
	}
	// Return messages as JSON
	response := map[string]interface{}{
		"type":     "group_messages",
		"messages": messages,
	}

	helper.RespondWithJSON(w, http.StatusOK, response)
}
