package handlers

import (
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func FollowHandler(w http.ResponseWriter, r *http.Request) {
	isFollowed := true

	w.Header().Set("Content-Type", "application/json")

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

	err = repository.Db.QueryRow("select privacy from  users  where  id  = ? ", id).Scan(&privacy)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	qCheck := `SELECT COUNT(*) FROM followers WHERE user_id = ? AND follower_id = ?`
	var count int
	err = repository.Db.QueryRow(qCheck, UserID, id).Scan(&count)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if privacy == "private" {
		// send a request message to the user
		return
	}
	if count > 0 {
		isFollowed = false
		_, err := repository.Db.Exec(`DELETE FROM followers WHERE  user_id = ? AND follower_id = ?`, UserID, id)
		if err != nil {
			http.Error(w, "Failed to follow user: "+err.Error(), http.StatusInternalServerError)
			return
		}

	} else {

		q := `INSERT INTO followers (user_id, follower_id) VALUES (?, ?)`
		_, err = repository.Db.Exec(q, UserID, id)
		if err != nil {
			http.Error(w, "Failed to follow user: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// w.WriteHeader(http.StatusOK)

	var followers int

	errr := repository.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE user_id = ?  `, id).Scan(&followers)
	if errr != nil {
		return
	}

	var following int

	errr = repository.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE follower_id = ?  `, id).Scan(&following)
	if errr != nil {
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"followers":  followers,
		"following":  following,
		"isFollowed": isFollowed,
	})
}
