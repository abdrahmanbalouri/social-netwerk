package model

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"social-network/internal/repository"
)

// InsertEvent saves the event in the database
type Eventt struct {
	ID          int64
	GroupID     string
	Title       string
	Description string
	Time        time.Time
}
type EVENT struct {
	ID          int    `json:"id"`
	GroupID     string `json:"group_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Time        string `json:"time"`
	CreatedAt   string `json:"created_at"`
}

func InsertEvent(db *sql.DB, e Eventt) (int64, error) {
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

func GetEvents(db *sql.DB, userID string, w http.ResponseWriter) ([]EVENT, error) {
	query := `
	SELECT e.id, e.group_id, e.title, e.description, e.time, e.created_at
	FROM events AS e
	JOIN event_Actions AS a ON e.id = a.event_id
	WHERE a.action = 'going' AND a.user_id = ?
	`

	rows, err := repository.Db.Query(query, userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return nil, err
	}
	defer rows.Close()
	var events []EVENT
	for rows.Next() {

		var event EVENT
		err := rows.Scan(
			&event.ID,
			&event.GroupID,
			&event.Title,
			&event.Description,
			&event.Time,
			&event.CreatedAt,
		)
		if err != nil {

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return nil, err
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Rows error: "+err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	return events, nil
}
