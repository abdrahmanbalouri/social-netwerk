package handlers

import (
	"encoding/json"
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
	"social-network/internal/repository/model"
)



func MyEavents(w http.ResponseWriter, r *http.Request) {
	userID, ok := helper.AuthenticateUser(r)
	if ok != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	events, err := model.GetEvents(repository.Db, userID, w)
	if err != nil {
		http.Error(w, "Failed to retrieve events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
