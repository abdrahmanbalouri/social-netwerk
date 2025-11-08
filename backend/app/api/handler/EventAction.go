package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"social-network/app/api/service"
	"social-network/app/helper"
)

func EventAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		helper.RespondWithError(w, http.StatusBadRequest, "invalid group ID")
		return
	}
	groupID := parts[4]

	var req struct {
		Action  string `json:"status"`
		EventID int    `json:"eventID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	err = service.HandleEventAction(userID, groupID, req.EventID, req.Action)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]any{
		"eventID": req.EventID,
		"status":  req.Action,
	})
}
