package handlers

import (
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func FollowRequest(w http.ResponseWriter, r *http.Request) {
	UserID, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	Fquery := `SELECT  u.id , u.first_name ,u.last_name, u.image
FROM users u 
join    follow_requests   f  on  f.follower_id = u.id  where user_id =   ? 
`
	rows, err := repository.Db.Query(Fquery, UserID)
	if err != nil {

		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var followRequest []map[string]interface{}
	for rows.Next() {
		var first_name ,last_name, profilePicture, idU string
		if err := rows.Scan(&idU, &first_name,&last_name, &profilePicture); err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		follower := map[string]interface{}{
			"id":       idU,
			"first_name": first_name,
			"last_name": last_name,
			"image":    profilePicture,
		}
		followRequest = append(followRequest, follower)
	}

	helper.RespondWithJSON(w, http.StatusOK, followRequest)
}
