package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

type GroupRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	InvitedUsers []string `json:"invitedUsers"`
}
type Group struct {
	ID          string
	Title       string
	Description string
}

func AddGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var newGroup GroupRequest
	if err := json.NewDecoder(r.Body).Decode(&newGroup); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if len(strings.TrimSpace(newGroup.Title)) == 0 || len(strings.TrimSpace(newGroup.Description)) == 0 {
		helper.RespondWithError(w, http.StatusBadRequest, "Title and description are required")
		return
	}

	// Get session cookie value.
	c, err := r.Cookie("session")
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "No valid session found")
		return
	}

	// Get the user's ID
	var adminID string
	if err := repository.Db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", c.Value).Scan(&adminID); err != nil {
		if err == sql.ErrNoRows {
			helper.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired session")
			return
		}
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user session")
		return
	}

	// Begin a transaction
	tx, err := repository.Db.Begin()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to start database transaction")
		return
	}
	defer tx.Rollback()

	grpID := helper.GenerateUUID()

	// Insert new group
	query1 := `INSERT INTO groups (id, title, description, admin_id) VALUES (?, ?, ?, ?)`
	if _, err := tx.Exec(query1, grpID, newGroup.Title, newGroup.Description, adminID); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create new group")
		return
	}

	// Insert the admin as a member of the group with is_admin set to true
	query2 := `INSERT INTO group_members (user_id, group_id) VALUES (?, ?)`
	if _, err := tx.Exec(query2, adminID, grpID); err != nil {
		fmt.Println("the error is : ", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to insert admin into group members table")
		return
	}

	// Process all invited users
	for _, nickname := range newGroup.InvitedUsers {
		var userID string

		err := tx.QueryRow(`SELECT id FROM users WHERE nickname = ?`, nickname).Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve invited user's ID")
			return
		}

		query3 := `INSERT INTO group_invitations (id, group_id, user_id, invited_by_user_id, status, created_at) VALUES (?, ?, ?, ?, ?, ?)`
		rowId := helper.GenerateUUID()
		createdAt := time.Now().UTC()
		if _, err := tx.Exec(query3, rowId, grpID, userID, adminID, "pending", createdAt); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to insert invited user into group_invitation table")
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	response := map[string]string{
		"message":  "Group and invitations created successfully",
		"group_id": grpID.String(),
	}
	helper.RespondWithJSON(w, http.StatusCreated, response)
}

func GetAllGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Check for user's session
	c, err := r.Cookie("session")
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "No valid session found")
		return
	}
	fmt.Println("session is :::", c)
	var userID string
	if err := repository.Db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", c.Value).Scan(&userID); err != nil {
		if err == sql.ErrNoRows {
			helper.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired session")
			return
		}
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user session")
		return
	}

	// get all groups
	query := `SELECT 
    g.id, 
    g.title, 
    g.description
	FROM groups g
	WHERE g.id NOT IN (
    SELECT gm.group_id 
    FROM group_members gm 
    WHERE gm.user_id = ? 
);`
	rows, err := repository.Db.Query(query, userID)
	if err != nil {
		fmt.Println("error is :", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "error getting all valid groups")
		return
	}

	var GroupJson []Group
	for rows.Next() {
		var g Group
		err := rows.Scan(&g.ID, &g.Title, &g.Description)
		if err != nil {
			fmt.Println("Failed to get group infos : ", err)
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get group infos")
			return
		}
		GroupJson = append(GroupJson, g)
	}

	// Return the posts as a JSON response
	helper.RespondWithJSON(w, http.StatusOK, GroupJson)
}

func GetMyGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Check for user's session
	c, err := r.Cookie("session")
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "No valid session found")
		return
	}
	fmt.Println("session is :::", c)
	var userID string
	if err := repository.Db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", c.Value).Scan(&userID); err != nil {
		if err == sql.ErrNoRows {
			helper.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired session")
			return
		}
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user session")
		return
	}

	// get all groups
	query := `SELECT 
    g.id, 
    g.title, 
    g.description
	FROM groups g
	WHERE g.id IN (
    SELECT gm.group_id 
    FROM group_members gm 
    WHERE gm.user_id = ? 
);`
	rows, err := repository.Db.Query(query, userID)
	if err != nil {
		fmt.Println("error is (my groups handler):", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "error getting all valid groups")
		return
	}

	var GroupJson []Group
	for rows.Next() {
		var g Group
		err := rows.Scan(&g.ID, &g.Title, &g.Description)
		if err != nil {
			fmt.Println("Failed to get group infos (my groups handler): ", err)
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get group infos")
			return
		}
		GroupJson = append(GroupJson, g)
	}

	// Return the posts as a JSON response
	helper.RespondWithJSON(w, http.StatusOK, GroupJson)
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
