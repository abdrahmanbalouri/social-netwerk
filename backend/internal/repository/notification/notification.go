package notification

import (
	"log"
	"net/http"

	"social-network/internal/repository"
)

type Notification struct {
	SenderID  string `json:"sender_id"`
	Name      string `json:"name"`
	Photo     string `json:"photo"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	CreatedAt int    `json:"created_at"`
}

func GetNotifications(id string, w http.ResponseWriter) ([]Notification, error) {
	rows, err := repository.Db.Query(`
	SELECT sender_id, type, message, created_at
	FROM notifications
	WHERE receiver_id = ?
	ORDER BY created_at DESC
	`, id)
	if err != nil {
		log.Println("Error fetching notifications:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification

	for rows.Next() {
		var notif Notification
		if err := rows.Scan(&notif.SenderID, &notif.Type, &notif.Message, &notif.CreatedAt); err != nil {
			log.Println("Error scanning notification:", err)
			continue
		}

		err = repository.Db.QueryRow(`
		SELECT nickname, image FROM users WHERE id = ?
		`, notif.SenderID).Scan(&notif.Name, &notif.Photo)
		if err != nil {
			log.Println("Error fetching sender info:", err)
			continue
		}

		notifications = append(notifications, notif)
	}
	return notifications, nil
}

func ClearNotifications(id string, w http.ResponseWriter) error {
	_, err := repository.Db.Exec(`
		DELETE FROM notifications WHERE receiver_id = ?
	`, id)
	if err != nil {
		log.Println("Error clearing notifications:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return err
	}
	return nil
}
