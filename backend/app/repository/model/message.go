package model

import (
	"social-network/app/repository"
)

type Message struct {
	Type           string `json:"type"`
	SubType        string `json:"subType"`
	SenderId       string `json:"senderId"`
	ReceiverId     string `json:"receiverId"`
	Content        string `json:"content"`
	MessageContent string `json:"messageContent"`
	CreatedAt      string `json:"createdAt"`
	First_name     string `json:"first_name"`
	Last_name      string `json:"last_name"`
	GroupID        string `json:"groupID"`
	Photo          string `json:"photo"`
}

func GetMessages(currentUserID, reciverId string) ([]Message, error) {
	query := `
		SELECT m.content, m.sender_id, m.receiver_id, m.sent_at , u.first_name,u.last_name,u.image FROM messages m
		LEFT JOIN users u ON u.id = m.sender_id
		WHERE (m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.sent_at DESC
		`
	rows, err := repository.Db.Query(query, currentUserID, reciverId, reciverId, currentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.Content, &msg.SenderId, &msg.ReceiverId, &msg.CreatedAt, &msg.First_name, &msg.Last_name, &msg.Photo); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func GetGroupMessages(currentUserID, groupId string) ([]Message, error) {
	err := repository.Db.QueryRow("SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?", groupId, currentUserID).Scan(new(any))
	if err != nil {
		return nil, err
	}

	query := `
		SELECT m.content, m.sender_id, m.sent_at FROM messages m
		WHERE m.group_id = ?
		ORDER BY m.sent_at DESC
		`

	rows, err := repository.Db.Query(query, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.Content, &msg.SenderId, &msg.CreatedAt); err != nil {
			return nil, err
		}
		err = repository.Db.QueryRow("SELECT first_name , last_name, image FROM users WHERE id = ?", msg.SenderId).Scan(&msg.First_name, &msg.Last_name, &msg.Photo)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}
