package handlers

import (
	"encoding/json"
	"net/http"

	"social-network/app/helper"
	"social-network/app/repository/notification"
)

// Notification structure to send to frontend

func Notifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, " method not allowed ")
		return
	}

	id, ok := helper.AuthenticateUser(r)
	if ok != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	seen := r.URL.Query().Get("bool")
	notifications, err := notification.GetNotifications(id, seen, w)
	if err != nil {
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
		return
	}

	w.WriteHeader(http.StatusOK)
}
