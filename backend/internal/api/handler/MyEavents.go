package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"social-network/internal/repository"
)

type Event struct {
	ID          int    `json:"id"`
	GroupID     int    `json:"group_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Time        string `json:"time"`
}

func MyEavents(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := sessionCookie.Value

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
		return
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.GroupID,
			&event.Title,
			&event.Description,
			&event.Time,
		)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		fmt.Println("11111111111")
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Rows error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
