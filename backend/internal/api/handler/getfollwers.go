package handlers

import (
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
	"social-network/internal/utils"
)

func Getfollowers(w http.ResponseWriter, r *http.Request) {
	UserId, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var users []struct {
		utils.User
	}
	Fquery := `SELECT  u.id , u.nickname, u.image
FROM followers f
JOIN users u ON u.id = f.user_id
WHERE f.follower_id = ?;
`
	rows, err := repository.Db.Query(Fquery, UserId)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Database query error")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user utils.User
		if err := rows.Scan(&user.ID, &user.Nickname, &user.Image); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Database scan error")
			return
		}
		users = append(users, struct {
			utils.User
		}{User: user})
	}

	if err := rows.Err(); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Database rows error")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, users)
}
