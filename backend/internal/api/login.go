package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"social-network/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type loginInformation struct {
	Nickname string `json:"email"`
	Password string `json:"password"`
}

func generateSessionID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db, err := repository.OpenDb()
	if err != nil {
		http.Error(w, "DB connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	var loginInformations loginInformation
	err = json.NewDecoder(r.Body).Decode(&loginInformations)
	fmt.Println(loginInformations)
	if err != nil {
	}

	var dbPassword string
	var userID int

	err = db.QueryRow("SELECT id, password FROM users WHERE( nickname = ? or email = ?)", loginInformations.Nickname, loginInformations.Nickname).Scan(&userID, &dbPassword)
	if err == sql.ErrNoRows || bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(loginInformations.Password)) != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionID := generateSessionID()
	_, err = db.Exec(`
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
