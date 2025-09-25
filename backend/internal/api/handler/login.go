package handlers

import (
	"encoding/json"
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
	"social-network/internal/repository/logindata"
	"social-network/internal/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginInformations utils.LoginInformation
	err := json.NewDecoder(r.Body).Decode(&loginInformations)
	if err != nil {
	}

	var dbPassword string
	var userID string

	err1 := logindata.Checklogindata(loginInformations.Nickname, repository.Db, w, &dbPassword, &userID, loginInformations.Password)
	if err1 == "0" {
		return
	}
	sessionID := helper.GenerateSessionID()
	_, err = repository.Db.Exec(`
		INSERT INTO sessions (user_id, token, expires_at)
		VALUES (?, ?, DATETIME('now', '+1 hour'))
	`, userID, sessionID)
	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	})
}
