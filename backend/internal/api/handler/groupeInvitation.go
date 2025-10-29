package handlers

import (
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func GroupeInvitation(w http.ResponseWriter, r *http.Request) {
	UserID, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	Fquery := `SELECT  g.id , g.title 
FROM groups g 
join    group_invitations   I  on  I.group_id = g.id  where I.user_id =   ? 
`
	rows, err := repository.Db.Query(Fquery, UserID)
	if err != nil {

		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var groupeInvitation []map[string]interface{}
	for rows.Next() {
		var title , idG string
		if err := rows.Scan(&idG, &title); err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		groups := map[string]interface{}{
			"id":       idG,
			"title": title,
	
		}
		groupeInvitation = append(groupeInvitation, groups)
	}

	helper.RespondWithJSON(w, http.StatusOK, groupeInvitation)
}
