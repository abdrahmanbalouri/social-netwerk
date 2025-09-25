package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"social-network/internal/helper"
)

func AddGroupHandler(w http.ResponseWriter, r *http.Request) {
	type GroupRequest struct {
		AdminID      string `json:"adminID"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if r.Method != http.MethodPost{
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var newGroup GroupRequest
	if err := json.NewDecoder(r.Body).Decode(&newGroup); err != nil {
		fmt.Println("111")
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if len(strings.TrimSpace(newGroup.AdminID)) == 0 || len(strings.TrimSpace(newGroup.Title)) == 0 || len(strings.TrimSpace(newGroup.Description)) == 0 {
		fmt.Println("2222")
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	query := `INSERT INTO groups (id, title, description, admin_id) VALUES (?, ?, ?, ?)`

	grpID := helper.GenerateUUID()
	fmt.Println("group id is :", grpID)
	
}
