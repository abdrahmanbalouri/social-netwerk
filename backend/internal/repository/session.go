package repository

import (
	"database/sql"
	"fmt"
	"time"
)

// Session represents a user session
type Session struct {
	UserID    string
	Token     string
	ExpiresAt time.Time
}

// GetSession retrieves a session by token
func GetSession(db *sql.DB, token string) (*Session, error) {
	row := db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE token=?", token)

	var userID string
	var expiresAtStr string
	err := row.Scan(&userID, &expiresAtStr)
	if err != nil {
		return nil, err
	}

	expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
	if err != nil {
		return nil, err
	}

	return &Session{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

// DeleteSession deletes a session by token
func DeleteSession(db *sql.DB, token string) error {
	mllm, err := db.Exec("DELETE FROM sessions WHERE token=?", token)
	fmt.Println(mllm)
	return err
}
