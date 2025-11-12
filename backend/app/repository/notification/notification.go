package notification

import (
	"log"
	"net/http"

	"social-network/pkg/db/sqlite"
)

type Notification struct {
	SenderID   string `json:"sender_id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Image      string `json:"image"`
	Type       string `json:"type"`
	Message    string `json:"message"`
	Seen       bool   `json:"seen"`
	CreatedAt  int    `json:"created_at"`
}

func GetNotifications(id, seen string, w http.ResponseWriter) ([]Notification, error) {
	if seen == "true" {
		_, err := sqlite.Db.Exec(`UPDATE notifications SET seen = TRUE WHERE seen = FALSE`)
		if err != nil {
		}

	}

	rows, err := sqlite.Db.Query(`
        SELECT n.sender_id,
               u.first_name,
               u.last_name,
               u.image,
               n.type,
               n.message,
               n.created_at,
			   n.seen
        FROM notifications n
        INNER JOIN users u ON u.id = n.sender_id
        WHERE n.receiver_id = ?
           OR EXISTS (
                SELECT 1
                FROM group_members gm
                WHERE gm.group_id = n.receiver_id
                  AND gm.user_id = ? and n.sender_id != gm.user_id
           )
        ORDER BY n.created_at DESC
    `, id, id)
	if err != nil {
		log.Println("Error fetching notifications:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification

	for rows.Next() {
		var notif Notification
		if err := rows.Scan(&notif.SenderID, &notif.First_name, &notif.Last_name, &notif.Image, &notif.Type, &notif.Message, &notif.CreatedAt, &notif.Seen); err != nil {
			log.Println("Error scanning notification:", err)
			continue
		}

		notifications = append(notifications, notif)
	}
	return notifications, nil
}

func ClearNotifications(id string, w http.ResponseWriter) error {
	_, err := sqlite.Db.Exec(`
		DELETE FROM notifications WHERE receiver_id = ?
	`, id)
	if err != nil {
		log.Println("Error clearing notifications:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return err
	}
	return nil
}
