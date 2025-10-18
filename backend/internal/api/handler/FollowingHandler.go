package handlers

import (
	"fmt"
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func FollowingHandler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("errrrrrrrrrrrrrrrrrrrrrrrrrrrrr", err)
		return
	}
	if privacy == "private" && isFollowing == 0 {

		helper.RespondWithJSON(w, http.StatusUnauthorized, "errrrrrrrrrrrrorrororororroororo")

		return
	}

	Fquery := `SELECT u.id ,  u.nickname,  u.image from followers f
	 JOIN users  u ON u.id = f.follower_id
	WHERE f.user_id = ?`

	rows, err := repository.Db.Query(Fquery, id)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var following []map[string]interface{}
	for rows.Next() {
		var username, profilePicture, idU string
		if err := rows.Scan(&idU, &username, &profilePicture); err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		follower := map[string]interface{}{
			"id":       idU,
			"nickname": username,
			"image":    profilePicture,
		}
		following = append(following, follower)
	}

	helper.RespondWithJSON(w, http.StatusOK, following)
}
