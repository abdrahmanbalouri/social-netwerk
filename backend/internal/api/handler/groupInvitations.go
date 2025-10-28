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
	InvitationID string `json:"invitation_id"`
	Response     string `json:"response"`
}
type GroupInvitation struct {
	// GroupID      string   `json:"groupID"`
	InvitedUsers []string `json:"invitedUsers"`
}

func GroupInvitationResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside group invitation response")
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

	var groupID string

	// i'll get the group id based on the invitatio id hit makantii9ch fl user :)
	err = repository.Db.QueryRow(`
        SELECT group_id FROM group_invitations 
        WHERE id = ? AND user_id = ?
    `, newResponse.InvitationID, userID).Scan(&groupID)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid invitation", http.StatusForbidden)
		return
	}
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if newResponse.Response == "accept" {
		query := `INSERT INTO group_members (user_id, group_id) VALUES (?, ?)`
		_, err = tx.Exec(query, userID, groupID)
		if err != nil {
			fmt.Println("error is :", err)
			helper.RespondWithError(w, http.StatusInternalServerError, "error inserting the user in the group member table")
			return
		}
	}
	query := `DELETE FROM group_invitations WHERE id = ?`
	_, err = tx.Exec(query, userID, newResponse.InvitationID)
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
	// fmt.Println("Body is :", json.NewDecoder(r.Body))
	if err := json.NewDecoder(r.Body).Decode(&newInvitation); err != nil {
		fmt.Println("Invalid request format : ", err)
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	fmt.Println("new invitation has :", newInvitation)
	// find the user id
	userID, IDerr := helper.AuthenticateUser(r)
	if IDerr != nil { //////////////////////////////////////////////////////
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
	invitationId := helper.GenerateUUID()

	for _, user := range newInvitation.InvitedUsers {
		// get the user's id
		/* 	var invitedUserID string
		invitedUserID = user */
		/* 	fmt.Println("user to invite is :", user)
		// bdelt     nicknamme bfirst naaaame  !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		query = `SELECT id FROM users WHERE nickname = ?`
		err = tx.QueryRow(query, user).Scan(&invitedUserID)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("11")
				continue
			}
			helper.RespondWithError(w, http.StatusInternalServerError, "Error finding the invited user")
			return
		} */
		// check if this user is already in the group or has a pending invit
		var exists1, exists2 bool
		query = `SELECT EXISTS (
			SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?
			UNION ALL
			SELECT 1 FROM group_invitations WHERE user_id = ? AND group_id = ?
			)`
		err = tx.QueryRow(query, user, GrpId, user, GrpId).Scan(&exists1)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Error checking for existing membership or invitation")
			return
		}

		fmt.Println("hna")
		if exists1 {
			continue
		} else {
			query = `SELECT EXISTS (
						SELECT 1
						FROM users
						WHERE id = $1
						);`
			err = tx.QueryRow(query, user).Scan(&exists2)
			if err != nil {
				fmt.Println("Database error is :", err)
				helper.RespondWithError(w, http.StatusInternalServerError, "Database error")
				return
			}
			if !exists2 {
				helper.RespondWithError(w, http.StatusBadRequest, "The invited user isn't a user of our website")
				return
			}
		}
		createdAt := time.Now().UTC()
		query = `INSERT INTO group_invitations (id, group_id, user_id, invited_by_user_id, status, created_at) VALUES (?, ?, ?, ?, ?, ?)`
		_, err = tx.Exec(query, invitationId, GrpId, user, userID, "pending", createdAt)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Error sending the invitation")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}
	response := map[string]any{
		"invitation_id": invitationId,
		"message":       "Invitation successfully processed",
	}
	fmt.Println("everything went good ----")
	helper.RespondWithJSON(w, http.StatusOK, response)
}
