package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	service "social-network/internal/api/service"
	"social-network/internal/utils"
	"strings"

	"social-network/internal/helper"
)

func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var newGroup utils.GroupRequest
	if err := json.NewDecoder(r.Body).Decode(&newGroup); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if len(strings.TrimSpace(newGroup.Title)) == 0 || len(strings.TrimSpace(newGroup.Description)) == 0 {
		helper.RespondWithError(w, http.StatusBadRequest, "Title and description are required")
		return
	}

	adminID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "No valid session found")
		return 
	}

	group, err := service.CreateNewGroup(adminID, newGroup)
	if err != nil {
		fmt.Println("Error creating group:", err)
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, group)
}

func GetAllGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		if err == sql.ErrNoRows {
			helper.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired session")
			return
		}
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user session")
		return
	}

	groups, err := service.GetAllAvailableGroups(userID)
	if err != nil {
		fmt.Println("error is :", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "error getting all valid groups")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, groups)
}

func GetMyGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	groups, err := service.GetUserGroups(userID)
	if err != nil {
		fmt.Println("error is (my groups handler):", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "error getting all valid groups")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, groups)
}

// {
//     "title": "group1",
//     "description": "description for group 1",
//     "invitedUsers": ["aya"]
// }

// {
//     "groupID": "group1",
//     "invitedUsers": "descritption for group 1"
// }
