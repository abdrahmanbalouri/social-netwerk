package logindata

import (
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Checklogindata(nickname string, db *sql.DB, w http.ResponseWriter, dbPassword *string, userID *string, Password string) string {
	err := db.QueryRow(
		"SELECT id, password FROM users WHERE  email = ? ",
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
