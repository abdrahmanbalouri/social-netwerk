package helper

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"html"
	"net/http"

	"social-network/internal/repository"
)

func GenerateSessionID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func Skip(str string) string {
	return html.EscapeString(str)
}

func AuthenticateUser(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", err
	}

	var userID string

	err = repository.Db.QueryRow(`
    SELECT u.id
    FROM sessions s
    JOIN users u ON s.user_id = u.id
    WHERE s.token = ?
`, cookie.Value).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}
