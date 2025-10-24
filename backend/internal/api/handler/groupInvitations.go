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

type GroupResponse struct {
	GrpID    string `json:"grpId"`
	Response string `json:"response"`
}
type GroupInvitation struct {
	// GroupID      string   `json:"groupID"`
	InvitedUsers []string `json:"invitedUsers"`
}

func GroupInvitationResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// check user's id using his session
	// check if the user is invited
	// check ht existance of the group
	// remove it from the group_invitations table
	// add this user to the group_members table

	var newResponse GroupResponse
	if err := json.NewDecoder(r.Body).Decode(&newResponse); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// find the user id
	userID, IDerr := helper.AuthenticateUser(r)
	if IDerr != nil {
		return
	}

	// start the transaction
	tx, err := repository.Db.Begin()
	if err != nil {
		fmt.Println("Failed to start database transaction")
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to start database transaction")
		return
	}
	defer tx.Rollback()

	// check if this user has an invitation
	var invitationID string
	if err := tx.QueryRow("SELECT id FROM group_invitations WHERE user_id = ? AND group_id = ?", userID, newResponse.GrpID).Scan(&invitationID); err != nil {
		if err == sql.ErrNoRows {
			helper.RespondWithError(w, http.StatusNotFound, "No pending invitation found for this user and group")
			return
		}
		fmt.Println("Failed to retrieve invitation")
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve invitation")
		return
	}

	if newResponse.Response == "accept" {
		query := `INSERT INTO group_members (user_id, group_id) VALUES (?, ?)`
		_, err = tx.Exec(query, userID, newResponse.GrpID)
		if err != nil {
			fmt.Println("error is :", err)
			helper.RespondWithError(w, http.StatusInternalServerError, "error inserting the user in the group member table")
			return
		}
	}
	query := `DELETE FROM group_invitations WHERE user_id = ? AND group_id = ?`
	_, err = tx.Exec(query, userID, newResponse.GrpID)
	if err != nil {
		fmt.Println("error deleting the invitation from it table")
		helper.RespondWithError(w, http.StatusInternalServerError, "error deleting the invitation from it table")
		return
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Failed to commit transaction")
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	response := map[string]string{
		"message": "Invitation successfully processed",
	}
	helper.RespondWithJSON(w, http.StatusOK, response)
}

func GroupInvitationRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusNotFound, "Group not found")
		return
	}
	GrpId := parts[3]

	var newInvitation GroupInvitation
	// if err := json.NewDecoder(r.Body).Decode(&newInvitation); err != nil {
	// 	helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
	// 	return
	// }

	// find the user id
	userID, IDerr := helper.AuthenticateUser(r)
	if IDerr != nil {
		return
	}

	// start the transaction
	tx, err := repository.Db.Begin()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to start database transaction")
		return
	}
	defer tx.Rollback()

	// check if the user is a member of that group
	var isMember bool
	query := `SELECT EXISTS (SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?)`
	if err := tx.QueryRow(query, userID, GrpId).Scan(&isMember); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to check group membership")
		return
	}

	if !isMember {
		helper.RespondWithError(w, http.StatusUnauthorized, "You are not a member of this group")
		return
	}

	for _, user := range newInvitation.InvitedUsers {
		// get the user's id
		var invitedUserID string
		query = `SELECT id FROM users WHERE nickname = ?`
		err = tx.QueryRow(query, user).Scan(&invitedUserID)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			helper.RespondWithError(w, http.StatusInternalServerError, "Error finding the invited user")
			return
		}
		// check if this user is already in the group or has a pending invit
		var exists bool
		query = `SELECT EXISTS (
        SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?
        UNION ALL
        SELECT 1 FROM group_invitations WHERE user_id = ? AND group_id = ?
    )`
		err = tx.QueryRow(query, invitedUserID, GrpId, invitedUserID, GrpId).Scan(&exists)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Error checking for existing membership or invitation")
			return
		}

		if exists {
			continue
		}
		rowId := helper.GenerateUUID()
		createdAt := time.Now().UTC()
		query = `INSERT INTO group_invitations (id, group_id, user_id, invited_by_user_id, status, created_at) VALUES (?, ?, ?, ?, ?, ?)`
		_, err = tx.Exec(query, rowId, GrpId, invitedUserID, userID, "pending", createdAt)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Error sending the invitation")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	response := map[string]string{
		"message": "Invitation successfully processed",
	}
	helper.RespondWithJSON(w, http.StatusOK, response)
}
