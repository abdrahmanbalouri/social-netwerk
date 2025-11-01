package handlers

import (
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func FollowersHandler(w http.ResponseWriter, r *http.Request) {
	UserID, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	var privacy string
	var isFollowing int

	err = repository.Db.QueryRow(`SELECT 
    u.privacy,
    CASE 
        WHEN f.follower_id IS NOT NULL THEN 1
        ELSE 0 
    END AS is_following
FROM users u
LEFT JOIN followers f 
    ON f.user_id = u.id      
   AND f.follower_id = ?      
WHERE u.id = ?;
`, UserID, id).Scan(&privacy, &isFollowing)
	if err != nil {
		return
	}
	if privacy == "private" && isFollowing == 0 && UserID != id {
		helper.RespondWithJSON(w, http.StatusUnauthorized, "errrrrrrrrrrrrorrororororroororo")
		return
	}

	Fquery := `SELECT  u.id , u.first_name, u.last_name, u.image
	FROM followers f
	JOIN users u ON u.id = f.follower_id
	WHERE f.user_id = ?;
	`
	rows, err := repository.Db.Query(Fquery, id)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var followers []map[string]interface{}
	for rows.Next() {
		var profilePicture, idU, first_name, last_name string
		if err := rows.Scan(&idU, &first_name, &last_name, &profilePicture); err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		follower := map[string]interface{}{
			"id":         idU,
			"first_name": first_name,
			"last_name":  last_name,
			"image":      profilePicture,
		}
		followers = append(followers, follower)
	}

	helper.RespondWithJSON(w, http.StatusOK, followers)
}
