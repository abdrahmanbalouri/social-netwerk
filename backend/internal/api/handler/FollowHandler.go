package handlers

import (
	"fmt"
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func FollowHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	UserID, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}
	qCheck := `SELECT COUNT(*) FROM followers WHERE user_id = ? AND follower_id = ?`
	var count int
	err = repository.Db.QueryRow(qCheck, UserID, id).Scan(&count)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Already following this user", http.StatusBadRequest)
		return
	}

	q := `INSERT INTO followers (user_id, follower_id) VALUES (?, ?)`
	_, err = repository.Db.Exec(q, UserID, id)
	if err != nil {
		fmt.Println("Error following user:", err)
		http.Error(w, "Failed to follow user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
