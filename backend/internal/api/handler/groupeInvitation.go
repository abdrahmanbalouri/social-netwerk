package handlers

import (
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func GroupeInvitation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	UserID, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	Fquery := `SELECT 
    g.id AS group_id,
    g.title,
    i.id AS invitation_id
	FROM groups g
	JOIN group_invitations i ON i.group_id = g.id
	WHERE i.user_id = ?;
	`
	rows, err := repository.Db.Query(Fquery, UserID)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var groupeInvitation []map[string]interface{}
	for rows.Next() {
		var title, idG, invitationID string
		if err := rows.Scan(&idG, &title, &invitationID); err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		groups := map[string]interface{}{
			"group_id":    idG,
			"title": title,
			"invitation_id": invitationID,
		}
		groupeInvitation = append(groupeInvitation, groups)
	}

	helper.RespondWithJSON(w, http.StatusOK, groupeInvitation)
}
