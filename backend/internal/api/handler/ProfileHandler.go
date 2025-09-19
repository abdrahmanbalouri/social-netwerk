package handlers

import (
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	q := `SELECT nickname, email, about, privacy, image FROM users WHERE id = ?`
	row := repository.Db.QueryRow(q, userID)

	var user struct {
		Nickname string
		Email    string
		About    string
		Privacy  string
		Image    string
	}

	if err := row.Scan(&user.Nickname, &user.Email, &user.About, &user.Privacy, &user.Image); err != nil {
		http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
		return
	}

	profileData := map[string]interface{}{
		"nickname": user.Nickname,
		"email":    user.Email,
		"about":    user.About,
		"privacy":  user.Privacy,
		"image":    user.Image,
	}

	helper.RespondWithJSON(w, http.StatusOK, profileData)
}
