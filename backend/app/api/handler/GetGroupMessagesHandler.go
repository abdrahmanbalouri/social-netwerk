package handlers

import (
	"database/sql"
	"net/http"

	"social-network/app/helper"
	"social-network/app/repository/model"
)

func GetGroupMessagesHandler(w http.ResponseWriter, r *http.Request) {
if r.Method !=  http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, " method not allowed ")
		return
}


	// Authenticate user

	currentUserID, err := helper.AuthenticateUser(r)
	if err != nil {

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
