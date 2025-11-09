package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"social-network/app/helper"
	"social-network/app/repository/model"
)

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// Authenticate user
	// w.Header().Set("Content-Type", "application/json")
	currentUserID, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Parse query parameters
	reciverId := r.URL.Query().Get("receiverId")
	if reciverId == "" {
		http.Error(w, "Missing receiverId parameter", http.StatusBadRequest)
		return
	}

	messages, err := model.GetMessages(currentUserID, reciverId)
	fmt.Println("rrrrrrrrrrrrrrrrrrrrrr", messages)
	if err == sql.ErrNoRows {
		fmt.Println(err)
		helper.RespondWithError(w, http.StatusNotFound, "No messages found")
		return
	} else if err != nil {
		fmt.Println(err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Error fetching messages: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"type":     "messages",
		"messages": messages,
	}

	helper.RespondWithJSON(w, http.StatusOK, response)
}
