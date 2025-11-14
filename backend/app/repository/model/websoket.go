package model

import (
	"database/sql"
	"fmt"
	"time"

	"social-network/pkg/db/sqlite"
)

// ✅ Jib l'ism w laqab dyal user b'id
func GetUserByID(currentUserID string) (map[string]any, error) {
	var first_name, last_name, photo, privacy string

	err := sqlite.Db.QueryRow(`SELECT first_name , last_name, image, privacy FROM users WHERE id = ?`, currentUserID).Scan(&first_name, &last_name, &photo, &privacy)
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
	err := sqlite.Db.QueryRow(`
				SELECT 1 FROM followers
				WHERE (user_id = ? AND follower_id = ?) OR (user_id = ? AND follower_id = ?)
			`, msg.ReceiverId, currentUserID, currentUserID, msg.ReceiverId).Scan(new(int))

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func CanSendMessageToReceiver(currentUserID string, msg Message) (bool, error) {
	err := sqlite.Db.QueryRow(`
				SELECT 1 FROM followers
				WHERE (user_id = ? AND follower_id = ?)
			`, currentUserID, msg.ReceiverId).Scan(new(int))
	if err == sql.ErrNoRows {
		var privacy string
		err := sqlite.Db.QueryRow(`
				SELECT privacy FROM users
				WHERE id = ?
			`, msg.ReceiverId).Scan(&privacy)
		if err != nil {
			return false, err
		}

		if privacy == "public" {
			return true, nil
		} else {
			return false, nil
		}
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// ✅ Insert notification direct sans check dyal follows
func SaveNotification(currentUserID string, msg Message) error {
	q := `INSERT INTO notifications ( sender_id, receiver_id, type, message, created_at) VALUES (?, ?, ?, ?, ?) `
	_, err := sqlite.Db.Exec(q, currentUserID, msg.ReceiverId, msg.Type, "Send you a message", time.Now())

	return err
}

func SaveMessage(currentUserID string, msg Message, imageFileName string) error {
	// Sauvegarder le message dans la base de données
	_, err := sqlite.Db.Exec(`
		INSERT INTO messages (sender_id, receiver_id, content, image)
		VALUES (?, ?, ?, ?)
	`, currentUserID, msg.ReceiverId, msg.MessageContent, imageFileName)
	if err != nil {
		return fmt.Errorf("failed to insert message: %v", err)
	}

	return nil
}

func IsUserGroupMember(currentUserID string, msg Message) (bool, error) {
	err := sqlite.Db.QueryRow("SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?", msg.GroupID, currentUserID).Scan(new(any))
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func SaveGroupMessage(currentUserID string, msg Message, imageFileName string) error {
	_, err := sqlite.Db.Exec(`
				INSERT INTO messages (group_id, sender_id, content, image)
				VALUES (?, ?, ?, ?)
			`, msg.GroupID, currentUserID, msg.MessageContent, imageFileName)

	return err
}

func GetGroupMembers(groupID string) ([]string, error) {
	var groupMembers []string
	rows, err := sqlite.Db.Query("SELECT user_id FROM group_members WHERE group_id = ?", groupID)
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

	err := sqlite.Db.QueryRow(`SELECT first_name , last_name, image FROM users WHERE id = ?`, currentUserID).Scan(&first_name, &last_name, &photo)
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

func IsFollowingRwauestReceiver(currentUserID string, msg Message) (bool, error) {
	var exist int
	err := sqlite.Db.QueryRow(`SELECT 1 FROM follow_requests WHERE user_id = ? AND follower_id = ?`, msg.ReceiverId, currentUserID).Scan(&exist)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func IsFollowingReceiver(currentUserID string, msg Message) (bool, error) {
	var exist int
	err := sqlite.Db.QueryRow(`SELECT 1 FROM followers WHERE user_id = ? AND follower_id = ?`, msg.ReceiverId, currentUserID).Scan(&exist)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func SaveFollowNotification(currentUserID string, msg Message) error {
	q := `INSERT INTO notifications ( sender_id, receiver_id, type, message, created_at) VALUES (?, ?, ?, ?, ?) `
	_, err := sqlite.Db.Exec(q, currentUserID, msg.ReceiverId, msg.Type, msg.MessageContent, time.Now().Unix())
	return err
}

func SaveGroupMessageNotification(currentUserID string, msg Message) error {
	groupm, err := GetGroupMembers(msg.GroupID)
	for _, v := range groupm {
		q := `INSERT INTO notifications ( sender_id, receiver_id, type, message, created_at) VALUES (?, ?, ?, ?, ?) `
		_, err = sqlite.Db.Exec(q, currentUserID, v, msg.Type, msg.MessageContent, time.Now().Unix())
	}
	return err
}

func SaveGroupInvitationNotification(currentUserID string, msg Message) error {
	q := `INSERT INTO notifications ( sender_id, receiver_id, type, message, created_at) VALUES (?, ?, ?, ?, ?) `
	_, err := sqlite.Db.Exec(q, currentUserID, msg.ReceiverId, msg.Type, msg.MessageContent, time.Now().Unix())
	return err
}

func Name(msg Message) (string, error) {
	var Name string
	err := sqlite.Db.QueryRow(`SELECT title FROM groups WHERE id = ?`, msg.GroupID).Scan(&Name)
	if err != nil {
		return "", err
	}
	return Name, nil
}

func SaveGroupJoinRequestNotification(currentUserID string, msg Message) (error, string) {
	var adminID string
	err := sqlite.Db.QueryRow(`SELECT admin_id FROM groups WHERE id = ?`, msg.ReceiverId).Scan(&adminID)
	if err != nil {
		return err, ""
	}
	_, err = sqlite.Db.Exec(`INSERT INTO notifications ( sender_id, receiver_id, type, message, created_at) VALUES (?, ?, ?, ?, ?) `, currentUserID, adminID, msg.Type, msg.MessageContent, time.Now().Unix())
	return err, adminID
}
