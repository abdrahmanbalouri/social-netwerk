package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func AddGroupHandler(w http.ResponseWriter, r *http.Request) {
	// session, err := r.Cookie("session")
	sessionID, err := r.Cookie("session_id")
	fmt.Println("session iD in groups handler is :", sessionID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"unauthorized"}`))
		return
	}

	type GroupRequest struct {
		AdminID     string `json:"adminID"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var newGroup GroupRequest
	if err := json.NewDecoder(r.Body).Decode(&newGroup); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if len(strings.TrimSpace(newGroup.AdminID)) == 0 || len(strings.TrimSpace(newGroup.Title)) == 0 || len(strings.TrimSpace(newGroup.Description)) == 0 {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// create new groups
	grpID := helper.GenerateUUID()
	query1 := `INSERT INTO groups (id, title, description, admin_id) VALUES (?, ?, ?, ?)`
	_, err1 := repository.Db.Exec(query1, grpID, newGroup.Title, newGroup.Description, newGroup.AdminID)
	if err1 != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create new group")
		return
	}

	// insert the admin in the grp_members table
	query2 := `INSERT INTO group_members (user_id, group_id) VALUES (?, ?)`
	_, err2 := repository.Db.Exec(query2, newGroup.AdminID, grpID)
	if err2 != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to insert the admin in the group_members table")
		return
	}
}
