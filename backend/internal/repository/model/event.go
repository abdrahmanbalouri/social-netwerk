package model

import (
	"database/sql"
	"time"
)

type Event struct {
	ID          int64
	GroupID     string
	Title       string
	Description string
	Time        time.Time
}

// InsertEvent saves the event in the database
func InsertEvent(db *sql.DB, e Event) (int64, error) {
	result, err := db.Exec(
		`INSERT INTO events (group_id, title, description, time) VALUES (?, ?, ?, ?)`,
		e.GroupID, e.Title, e.Description, e.Time.Format(time.RFC3339),
	)
	if err != nil {
		return 0, err
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return eventID, nil
}

// CheckUserInGroup verifies if a user belongs to a group
func CheckUserInGroup(db *sql.DB, userID, groupID string) (bool, error) {
	var exists string
	err := db.QueryRow(`SELECT user_id FROM group_members WHERE user_id=? AND group_id=?`, userID, groupID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
