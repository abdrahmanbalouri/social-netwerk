package model

import (
	"database/sql"
	"time"

	"social-network/internal/repository"
)

// ✅ Jib l'ism w laqab dyal user b'id
func GetUserByID(currentUserID string) (map[string]any, error) {
	var first_name, last_name, photo, privacy string

	err := repository.Db.QueryRow(`SELECT first_name , last_name, image, privacy FROM users WHERE id = ?`, currentUserID).Scan(&first_name, &last_name, &photo, &privacy)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	user := map[string]any{
		"first_name": first_name,
		"last_name":  last_name,
		"photo":      photo,
		"privacy":    privacy,
	}
	return user, nil
}

// ✅ Vérifie ghir wach receiver kayn
func CheckIfUsersFollowEachOther(currentUserID string, msg Message) (bool, error) {
	var exist int
	err := repository.Db.QueryRow(`
				SELECT 1 FROM followers
				WHERE (user_id = ? AND follower_id = ?) OR (user_id = ? AND follower_id = ?)
			`, currentUserID, msg.ReceiverId, msg.ReceiverId, currentUserID).Scan(&exist)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// ✅ Insert notification direct sans check dyal follows
func SaveNotification(currentUserID string, msg Message) error {
	q := `INSERT INTO notifications ( sender_id, receiver_id, type, message, created_at) VALUES (?, ?, ?, ?, ?) `
	_, err := repository.Db.Exec(q, currentUserID, msg.ReceiverId, msg.Type, "Send you a message", time.Now().Unix())

	return err
}

// ✅ Insert message direct sans check dyal follows
func SaveMessage(currentUserID string, msg Message) error {
	_, err := repository.Db.Exec(`
				INSERT INTO messages (sender_id, receiver_id, content)
				VALUES (?, ?, ?)
			`, currentUserID, msg.ReceiverId, msg.MessageContent)
	return err
}

func IsUserGroupMember(currentUserID string, msg Message) (bool, error) {
	err := repository.Db.QueryRow("SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?", msg.GroupID, currentUserID).Scan(new(any))
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func SaveGroupMessage(currentUserID string, msg Message) error {
	_, err := repository.Db.Exec(`
				INSERT INTO messages (group_id, sender_id, content)
				VALUES (?, ?, ?)
			`, msg.GroupID, currentUserID, msg.MessageContent)

	return err
}

func GetGroupMembers(groupID string) ([]string, error) {
	var groupMembers []string
	rows, err := repository.Db.Query("SELECT user_id FROM group_members WHERE group_id = ?", groupID)
	if err != nil {
		return groupMembers, err
	}
	defer rows.Close()

	var userID string
	for rows.Next() {
		if err := rows.Scan(&userID); err != nil {
			return groupMembers, err
		}
		groupMembers = append(groupMembers, userID)
	}
	return groupMembers, nil
}

func GetUserInfoByID(currentUserID string) (map[string]any, error) {
	var first_name, last_name, photo string

	err := repository.Db.QueryRow(`SELECT first_name , last_name, image FROM users WHERE id = ?`, currentUserID).Scan(&first_name, &last_name, &photo)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	user := map[string]any{
		"first_name": first_name,
		"last_name":  last_name,
		"photo":      photo,
	}
	return user, nil
}

func IsFollowingReceiver(currentUserID string, msg Message) (bool, error) {
	var exist int
	err := repository.Db.QueryRow(`SELECT 1 FROM followers WHERE user_id = ? AND follower_id = ?`, msg.ReceiverId, currentUserID).Scan(&exist)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func SaveGroupInvitationNotification(currentUserID string, msg Message) error {
	q := `INSERT INTO notifications ( sender_id, receiver_id, type, message, created_at) VALUES (?, ?, ?, ?, ?) `
	_, err := repository.Db.Exec(q, currentUserID, msg.ReceiverId, msg.Type, msg.MessageContent, time.Now().Unix())
	return err
}
