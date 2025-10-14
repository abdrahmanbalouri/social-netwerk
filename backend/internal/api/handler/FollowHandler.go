package handlers

import (
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func FollowHandler(w http.ResponseWriter, r *http.Request) {
	isFollowed := true
	isPending := false

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
	err = repository.Db.QueryRow("SELECT privacy FROM users WHERE id = ?", id).Scan(&privacy)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	qCheck := `SELECT COUNT(*) FROM followers WHERE user_id = ? AND follower_id = ?`
	var count int
	err = repository.Db.QueryRow(qCheck, id, UserID).Scan(&count)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if privacy == "private" {
		var reqCount int
		repository.Db.QueryRow("SELECT COUNT(*) FROM follow_requests WHERE user_id = ? AND follower_id = ?", id, UserID).Scan(&reqCount)

		if reqCount > 0 {
			_, err = repository.Db.Exec("DELETE FROM follow_requests WHERE user_id = ? AND follower_id = ?", id, UserID)
			if err != nil {
				http.Error(w, "Failed to cancel follow request: "+err.Error(), http.StatusInternalServerError)
				return
			}
			isPending = false
		} else {
			_, err = repository.Db.Exec("INSERT INTO follow_requests (user_id, follower_id) VALUES (?, ?)", id, UserID)
			if err != nil {
				http.Error(w, "Failed to create follow request: "+err.Error(), http.StatusInternalServerError)
				return
			}
			isPending = true
		}

		var followers, following int
		repository.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE user_id = ?`, id).Scan(&followers)
		repository.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE follower_id = ?`, id).Scan(&following)
		helper.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"followers":  followers,
			"following":  following,
			"isFollowed": false,
			"isPending":  isPending,
		})
		return
	}

	if count > 0 {
		isFollowed = false
		_, err := repository.Db.Exec(`DELETE FROM followers WHERE user_id = ? AND follower_id = ?`, id, UserID)
		if err != nil {
			http.Error(w, "Failed to unfollow: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		q := `INSERT INTO followers (user_id, follower_id) VALUES (?, ?)`
		_, err = repository.Db.Exec(q, id, UserID)
		if err != nil {
			http.Error(w, "Failed to follow user: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var followers, following int
	repository.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE user_id = ?`, id).Scan(&followers)
	repository.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE follower_id = ?`, id).Scan(&following)

	helper.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"followers":  followers,
		"following":  following,
		"isFollowed": isFollowed,
		"isPending":  false,
	})
}
