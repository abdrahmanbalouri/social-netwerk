package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	ReceiverId     string `json:"receiverId"`
	MessageContent string `json:"messageContent"`
	Type           string `json:"type"`
}

var ConnectedUsers = make(map[string][]*websocket.Conn)

func Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	id, _ := helper.AuthenticateUser(r)
	ConnectedUsers[id] = append(ConnectedUsers[id], conn)
	defer conn.Close()
	for {
		_, msg, err := conn.NextReader()
		if err != nil {
			fmt.Println("Connection closed:", err)
			break
		}

		var messageStruct Message
		decoder := json.NewDecoder(msg)
		err = decoder.Decode(&messageStruct)
		if err != nil {
			fmt.Println("err", err)
			continue
		}

		if messageStruct.Type == "follow" {
			q := `SELECT id FROM followers WHERE user_id = ? AND follower_id = ?`
			var followID int

			_ = repository.Db.QueryRow(q, id, messageStruct.ReceiverId).Scan(&followID)

			if followID != 0 {
				messageStruct.Type="already_following"
			}
			messageStruct.ReceiverId = id
		}

		for i, conArr := range ConnectedUsers {
			if i != id {
				for _, con := range conArr {
					jsonMsg, err := json.Marshal(messageStruct)
					if err != nil {
						fmt.Println("Error marshaling message:", err)
						continue
					}
					if err := con.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
						fmt.Println("Error writing message:", err)
						break
					}
				}
			}
		}
	}
}
