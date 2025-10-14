package handlers

import (
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
	"social-network/internal/utils"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}
	rows, err := repository.Db.Query("SELECT id, nickname, image  FROM users")
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	defer rows.Close()

	var users []struct {
		utils.User
	}
	for rows.Next() {
		var user struct {
			utils.User
		}
		if err := rows.Scan(&user.ID, &user.Nickname, &user.Image); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to process users")
			return
		}
		if user.ID == userID {
			continue
		}
		users = append(users, user)
	}

	helper.RespondWithJSON(w, http.StatusOK, users)
}
