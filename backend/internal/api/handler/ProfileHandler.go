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

	userID := r.URL.Query().Get("userId")
	if userID == "0" {
		userid, err := helper.AuthenticateUser(r)
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		userID = userid

	}
	q := `SELECT id, nickname, email, about, privacy, image, cover FROM users WHERE id = ?`
	row := repository.Db.QueryRow(q, userID)

	var user struct {
		id       string
		Nickname string
		Email    string
		About    string
		Privacy  string
		Image    string
		Cover    string
	}

	if err := row.Scan(&user.id, &user.Nickname, &user.Email, &user.About, &user.Privacy, &user.Image, &user.Cover); err != nil {
		http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
		return
	}

	profileData := map[string]interface{}{
		"id":       user.id,
		"nickname": user.Nickname,
		"email":    user.Email,
		"about":    user.About,
		"privacy":  user.Privacy,
		"image":    user.Image,
		"cover":    user.Cover,
	}

	helper.RespondWithJSON(w, http.StatusOK, profileData)
}
