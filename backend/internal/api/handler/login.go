package handlers

import (
	"encoding/json"
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
	"social-network/internal/repository/login"
	"social-network/internal/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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
		repository.Db,
		w,
		&dbPassword,
		&userID,
		loginData.Password,
	)
	if err1 != "" {
		http.Error(w, err1, http.StatusUnauthorized)
		return
	}

	// Generate session ID
	sessionID := helper.GenerateSessionID()

	// Store session in database
	err = logindata.CreateSession(userID, sessionID)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Login successful"}`))
}
