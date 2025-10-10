package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	Name           string `json:"name"`
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
			q := `
				SELECT u.nickname 
				FROM followers f
				JOIN users u ON u.id = f.user_id
				WHERE f.user_id = ? AND f.follower_id = ?
			`

			var nickname string
			err := repository.Db.QueryRow(q, id, messageStruct.ReceiverId).Scan(&nickname)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					messageStruct.Type = "follow"
					messageStruct.Name = nickname
					fmt.Println(nickname)
				} else {
					log.Println("Error checking follow:", err)
				}
			} else {
				messageStruct.Type = "already_following"
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
