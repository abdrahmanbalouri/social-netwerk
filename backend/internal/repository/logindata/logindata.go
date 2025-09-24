package logindata

import (
	"database/sql"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Checklogindata(nickname string, db *sql.DB, w http.ResponseWriter, dbPassword *string, userID *string, Password string) string {
    err := db.QueryRow(
        "SELECT id, password FROM users WHERE (nickname = ? OR email = ?)",
        nickname, nickname,
    ).Scan(userID, dbPassword)
    
    if err == sql.ErrNoRows || bcrypt.CompareHashAndPassword([]byte(*dbPassword), []byte(Password)) != nil {
        fmt.Println("no ", err)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusUnauthorized)
        return "0"
    } else if err != nil {
        fmt.Println("no")

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        return "0"
    }

    return "1"
}

