package model

import "social-network/app/repository"

func GetEventByIDAndGroup(eventID int, groupID string) (int, string, error) {
	var id int
	var time string
	err := repository.Db.QueryRow(`SELECT id, time FROM events WHERE id = ? AND group_id = ?`, eventID, groupID).Scan(&id, &time)
	return id, time, err
}

func GetUserEventAction(eventID int, userID string) (string, error) {
	var action string
	err := repository.Db.QueryRow(`SELECT action FROM event_actions WHERE event_id = ? AND user_id = ?`, eventID, userID).Scan(&action)
	return action, err
}

func InsertEventAction(eventID int, userID string, action string) error {
	_, err := repository.Db.Exec(`INSERT INTO event_actions (event_id, user_id, action) VALUES (?, ?, ?)`, eventID, userID, action)
	return err
}

func UpdateEventAction(eventID int, userID string, action string) error {
	_, err := repository.Db.Exec(`UPDATE event_actions SET action = ? WHERE event_id = ? AND user_id = ?`, action, eventID, userID)
	return err
}

func DeleteEventAction(eventID int, userID string) error {
	_, err := repository.Db.Exec(`DELETE FROM event_actions WHERE event_id = ? AND user_id = ?`, eventID, userID)
	return err
}
