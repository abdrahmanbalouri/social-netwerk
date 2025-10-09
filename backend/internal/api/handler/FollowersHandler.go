package handlers

import (
	"fmt"
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func FollowersHandler(w http.ResponseWriter, r *http.Request) {
	/* UserID, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	} */
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	Fquery := `SELECT u.username,  u.profile_picture from users u
	INNER JOIN followers f ON u.id = f.follower_id
	WHERE f.follower_id = ?`
	rows, err := repository.Db.Query(Fquery, id)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var followers []map[string]interface{}
	for rows.Next() {
		var username, profilePicture string
		if err := rows.Scan(&username, &profilePicture); err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		follower := map[string]interface{}{
			"username":        username,
			"profile_picture": profilePicture,
		}
		followers = append(followers, follower)
	}

	helper.RespondWithJSON(w, http.StatusOK, followers)

	fmt.Println("folowers", followers)
}
