package handlers

import (
	"encoding/json"
	"net/http"

	"social-network/app/helper"
	"social-network/app/repository/model"
	"social-network/pkg/db/sqlite"
)

func MyEavents(w http.ResponseWriter, r *http.Request) {
	userID, ok := helper.AuthenticateUser(r)
	if ok != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	events, err := model.GetEvents(sqlite.Db, userID, w)
	if err != nil {
		http.Error(w, "Failed to retrieve events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
