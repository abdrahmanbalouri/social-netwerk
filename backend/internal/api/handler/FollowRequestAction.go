package handlers

import (
	"encoding/json"
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func FollowRequestAction(w http.ResponseWriter, r *http.Request) {
	UserID, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	var req struct {
		ID     string    `json:"id"`    
		Action string `json:"action"` 
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM follow_requests WHERE user_id = ? AND follower_id = ?)`
	err = repository.Db.QueryRow(checkQuery, UserID, req.ID).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Follow request not found", http.StatusNotFound)
		return
	}

	if req.Action == "accept" {
		_, err = repository.Db.Exec(
			"INSERT INTO followers (user_id, follower_id) VALUES (?, ?)",
			UserID, req.ID,
		)
		if err != nil {
			http.Error(w, "Error inserting follower: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	_, err = repository.Db.Exec(
		"DELETE FROM follow_requests WHERE user_id = ? AND follower_id = ?",
		UserID, req.ID,
	)
	if err != nil {
		http.Error(w, "Error deleting follow request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]string{
		"status": "success",
		"action": req.Action,
	})
}
