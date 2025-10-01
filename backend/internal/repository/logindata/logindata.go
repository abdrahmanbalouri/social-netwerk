package logindata

import (
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Checklogindata(nickname string, db *sql.DB, w http.ResponseWriter, dbPassword *string, userID *string, Password string) string {
	err := db.QueryRow(
		"SELECT id, password FROM users WHERE (nickname = ? OR email = ?)",
		nickname, nickname,
	).Scan(userID, dbPassword)

	if err == sql.ErrNoRows || bcrypt.CompareHashAndPassword([]byte(*dbPassword), []byte(Password)) != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return "0"
	} else if err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return "0"
	}

	return "1"
}
