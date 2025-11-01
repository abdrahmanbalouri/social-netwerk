package handlers

import (
	"fmt"
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
	"social-network/internal/utils"
)

func GetCommunFriends(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("GetCommunFriends called")
	if r.Method != "GET" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
	
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	rows, err := repository.Db.Query(`
	SELECT DISTINCT u.id, u.first_name ,u.last_name, u.image
	FROM users u
	INNER JOIN followers f 
	ON (u.id = f.follower_id OR u.id = f.user_id)
	WHERE (f.user_id = ? OR f.follower_id = ?);
	`, userID, userID)
	if err != nil {
		fmt.Println("Error querying common friends:", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch common friends")
		return
	}
	defer rows.Close()

	var friends []struct {
		utils.User
	}
	for rows.Next() {
		var friend struct {
			utils.User
		}
		if err := rows.Scan(&friend.ID, &friend.First_name, &friend.Last_name, &friend.Image); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to process common friends")
			return
		}
		if userID != friend.ID {
			friends = append(friends, friend)
		}
	}

	helper.RespondWithJSON(w, http.StatusOK, friends)
}
