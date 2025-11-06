package handlers

import (
	"encoding/json"
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository/notification"
)

// Notification structure to send to frontend

func Notifications(w http.ResponseWriter, r *http.Request) {
	id, ok := helper.AuthenticateUser(r)
	if ok != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	notifications, err := notification.GetNotifications(id, w)
	if err != nil {
		http.Error(w, "Failed to get notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func ClearNotifications(w http.ResponseWriter, r *http.Request) {
	id, ok := helper.AuthenticateUser(r)
	if ok != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	err := notification.ClearNotifications(id, w)
	if err != nil {
		http.Error(w, "Failed to clear notifications", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
