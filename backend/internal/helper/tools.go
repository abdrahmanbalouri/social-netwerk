package helper

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"social-network/internal/repository"

	"github.com/gofrs/uuid/v5"
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
	fmt.Println(payload)
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
	// ("Cookie is :", cookie)

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

func GenerateUUID() uuid.UUID {
	// Create a Version 4 UUID
	u2, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	return u2
}
