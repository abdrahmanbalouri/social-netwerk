package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func GetGroupMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// Authenticate user

	_, err := helper.AuthenticateUser(r)
	if err != nil {
		fmt.Println("Authentication error:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Parse query parameters
	groupId := r.URL.Query().Get("groupId")
	if groupId == "" {
		http.Error(w, "Missing receiverId parameter", http.StatusBadRequest)
		return
	}
	query := `
		SELECT m.content, m.sender_id, m.sent_at FROM messages m
		WHERE m.group_id = ?
		ORDER BY m.sent_at DESC
		`

	rows, err := repository.Db.Query(query, groupId)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	type Message struct {
		Content    string `json:"content"`
		SenderId   string `json:"senderId"`
		CreatedAt  string `json:"createdAt"`
		First_name string `json:"first_name"`
		Last_name  string `json:"last_name"`
	}

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.Content, &msg.SenderId, &msg.CreatedAt); err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		err = repository.Db.QueryRow("SELECT first_name , last_name FROM users WHERE id = ?", msg.SenderId).Scan(&msg.First_name, &msg.Last_name)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Return messages as JSON
	response := map[string]interface{}{
		"type":     "group_messages",
		"messages": messages,
	}

	json.NewEncoder(w).Encode(response)
}
