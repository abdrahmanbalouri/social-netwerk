package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
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
		fmt.Println("Authentication error:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Parse query parameters
	reciverId := r.URL.Query().Get("receiverId")
	if reciverId == "" {
		http.Error(w, "Missing receiverId parameter", http.StatusBadRequest)
		return
	}
	query := `
		SELECT m.content, m.sender_id, m.receiver_id, m.sent_at FROM messages m
		WHERE (m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.sent_at DESC
		`

	rows, err := repository.Db.Query(query, currentUserID, reciverId, reciverId, currentUserID)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	type Message struct {
		Content    string `json:"content"`
		SenderId   string `json:"senderId"`
		ReceiverId string `json:"receiverId"`
		CreatedAt  string `json:"createdAt"`
		Username   string `json:"username"`
	}

	var messages []Message
	for rows.Next() {
		var msg Message
		err = repository.Db.QueryRow("SELECT nickname FROM users WHERE id = ?", currentUserID).Scan(&msg.Username)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if err := rows.Scan(&msg.Content, &msg.SenderId, &msg.ReceiverId, &msg.CreatedAt); err != nil {
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
		"type":     "messages",
		"messages": messages,
	}

	json.NewEncoder(w).Encode(response)
}
