package logindata

import (
	"database/sql"
	"net/http"

	"social-network/app/repository"

	"golang.org/x/crypto/bcrypt"
)

func Checklogindata(nickname string, db *sql.DB, w http.ResponseWriter, dbPassword *string, userID *string, Password string) string {
	if nickname == "" || Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return "email or password incorrect"
	}
	err := db.QueryRow(
		"SELECT id, password FROM users WHERE  ( nickname = ? or email = ?)",
		nickname, nickname,
	).Scan(userID, dbPassword)

	if err == sql.ErrNoRows || bcrypt.CompareHashAndPassword([]byte(*dbPassword), []byte(Password)) != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return "email or password incorrect"
	} else if err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return "database error"
	}

	return ""
}

func CreateSession(userID, token string) error {
	_, err := repository.Db.Exec(`
		INSERT INTO sessions (user_id, token, expires_at)
		VALUES (?, ?, DATETIME('now', '+1 hour'))
	`, userID, token)
	return err
}
