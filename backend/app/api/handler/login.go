package handlers

import (
	"encoding/json"
	"net/http"

	"social-network/app/helper"
	logindata "social-network/app/repository/login"
	"social-network/app/utils"
	"social-network/pkg/db/sqlite"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	var loginData utils.LoginInformation
	// Parse JSON body
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Use existing Checklogindata for validation
	var dbPassword, userID string
	err1 := logindata.Checklogindata(
		loginData.Nickname,
		sqlite.Db,
		w,
		&dbPassword,
		&userID,
		loginData.Password,
	)
	if err1 != "" {
		helper.RespondWithError(w, http.StatusUnauthorized, err1)
		return
	}

	// Generate session ID
	sessionID := helper.GenerateSessionID()

	// Store session in database
	err = logindata.CreateSession(userID, sessionID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	})

	// Respond with success
	helper.RespondWithJSON(w, http.StatusOK, "Login successful")
}
