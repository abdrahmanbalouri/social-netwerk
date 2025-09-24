package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"social-network/internal/helper"
)

func AddGroupHandler(w http.ResponseWriter, r *http.Request) {
	type GroupRequest struct {
		UserID      string `json:"userid"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	var newGroup GroupRequest
	if err := json.NewDecoder(r.Body).Decode(&newGroup); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if len(strings.TrimSpace(newGroup.UserID)) == 0 || len(strings.TrimSpace(newGroup.Title)) == 0 || len(strings.TrimSpace(newGroup.Description)) == 0 {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	query := `INSERT INTO groups (id, title, description, admin_id) VALUES (?, ?, ?, ?)`
	
}
