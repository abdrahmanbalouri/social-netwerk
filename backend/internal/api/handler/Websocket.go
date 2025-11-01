package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"sync"
	"time"

	"social-network/internal/helper"
	"social-network/internal/repository"

	"github.com/gorilla/websocket"
)

// Global WebSocket-related variables managed within this package
var (
	Clients      = make(map[string][]*websocket.Conn)
	ClientsMutex sync.Mutex
	Upgrader     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow connections from localhost:3000
			return r.Header.Get("Origin") == "http://localhost:3000"
		},
	}
)

type Message struct {
	Type           string `json:"type"`
	SubType        string `json:"subType"`
	ReceiverId     string `json:"receiverId"`
	MessageContent string `json:"messageContent"`
	First_name     string `json:"first_name"`
	Last_name      string `json:"last_name"`
	GroupID        string `json:"groupID"`
	Photo          string `json:"photo"`
}

func GetOnlineUsers() []string {
	users := []string{}
	for id := range Clients {
		users = append(users, id)
	}
	return users
}

func BrodcastOnlineListe() {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()
	msg := map[string]any{
		"type":  "online_list",
		"users": GetOnlineUsers(),
	}
	for _, client := range Clients {
		for _, conn := range client {
			if err := conn.WriteJSON(msg); err != nil {
				log.Println("WebSocket write error:", err)
				if err := conn.Close(); err != nil {
					log.Println("WebSocket close error:", err)
				}
			}
		}
	}
}

// WebSocketHandler handles WebSocket connections and manages user sessions
func Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	currentUserID, err := helper.AuthenticateUser(r)
	if err != nil {
		log.Println("Authentication error:", err)
		return
	}

	ClientsMutex.Lock()
	Clients[currentUserID] = append(Clients[currentUserID], conn)
	ClientsMutex.Unlock()

	BrodcastOnlineListe()
	// Listen for incoming messages
	Loop(conn, currentUserID)

	defer func() {
		ClientsMutex.Lock()
		conns, ok := Clients[currentUserID]

		if ok {
			for i, c := range conns {
				if c == conn {
					Clients[currentUserID] = append(conns[:i], conns[i+1:]...)
					break
				}
			}
			if len(Clients[currentUserID]) == 0 {
				delete(Clients, currentUserID)
				BrodcastOnlineStatus(currentUserID, false)
				if err := conn.Close(); err != nil {
					log.Println("WebSocket close error:", err)
				}
			}
		}
		ClientsMutex.Unlock()
	}()
}

func Loop(conn *websocket.Conn, currentUserID string) {
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		var first_name, last_name string
		err := repository.Db.QueryRow(`SELECT first_name , last_name FROM users WHERE id = ?`, currentUserID).Scan(&first_name, &last_name)
		if err != nil {
			log.Println("DB error getting sender firsy_name && last_name:", err)
			first_name = "Unknown"
			last_name = "Unknown"

		}
		switch msg.Type {
		case "logout":
			BrodcastOnlineStatus(currentUserID, false)
		//  HANDLE CHAT MESSAGE
		case "online_list":
			// ✅ Send current online list to the new user
			BrodcastOnlineListe()
		case "message":

			if msg.ReceiverId == "" {
				log.Println("Invalid recipient ID")
				continue
			}

			// ✅ Vérifie ghir wach receiver kayn
			var exist int
			err := repository.Db.QueryRow(`
				SELECT 1 FROM followers
				WHERE (user_id = ? AND follower_id = ?) OR (user_id = ? AND follower_id = ?)
			`, currentUserID, msg.ReceiverId, msg.ReceiverId, currentUserID).Scan(&exist)

			if err == sql.ErrNoRows {
				log.Println("Recipient user not found")
				continue
			} else if err != nil {
				log.Println("DB error:", err)
				continue
			}
			if exist == 0 {
				log.Println("Users are not connected as followers")
				continue
			}

			q := `INSERT INTO notifications ( sender_id, receiver_id, type, message, created_at) VALUES (?, ?, ?, ?, ?) `
			_, _ = repository.Db.Exec(q, currentUserID, msg.ReceiverId, msg.Type, "Send you a message", time.Now().Unix())

			// ✅ Insert message direct sans check dyal follows
			_, err = repository.Db.Exec(`
				INSERT INTO messages (sender_id, receiver_id, content)
				VALUES (?, ?, ?)
			`, currentUserID, msg.ReceiverId, msg.MessageContent)
			if err != nil {
				log.Println("DB error inserting message:", err)
				continue
			}

			// ✅ send message to receiver
			sendToUser(msg.ReceiverId, map[string]any{
				"type":    msg.Type,
				"from":    currentUserID,
				"to":      msg.ReceiverId,
				"content": msg.MessageContent,
				"time":    time.Now().Format(time.RFC3339),
			})

			// ✅ send back to sender also
			sendToUser(currentUserID, map[string]any{
				"type":    msg.Type,
				"from":    currentUserID,
				"to":      msg.ReceiverId,
				"content": msg.MessageContent,
				"time":    time.Now().Format(time.RFC3339),
			})

			// send notification to receiver
			BrodcastNotification(msg.ReceiverId, map[string]any{
				"type":    "notification",
				"subType": "message",
				"from":    currentUserID,
				"content": "sent you a message",
				"name":    first_name + " " + last_name,
				"time":    time.Now().Format(time.RFC3339),
			})

		// ===============================
		//  HANDLE FOLLOW NOTIFICATION
		// ===============================
		case "group_message":
			if msg.GroupID == "" {
				log.Println("Invalid group ID")
				continue
			}
			err := repository.Db.QueryRow("SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?", msg.GroupID, currentUserID).Scan(new(interface{}))
			if err == sql.ErrNoRows {
				log.Println("User is not a member of the group")
				continue
			} else if err != nil {
				log.Println("DB error:", err)
				continue
			}


			_, err = repository.Db.Exec(`
				INSERT INTO messages (group_id, sender_id, content)
				VALUES (?, ?, ?)
			`, msg.GroupID, currentUserID, msg.MessageContent)
			if err != nil {
				log.Println("DB error inserting group message:", err)
				continue
			}

			sendToGroupMembers(msg.GroupID, currentUserID, map[string]any{
				"type":     msg.Type,
				"from":     currentUserID,
				"groupID":  msg.GroupID,
				"content":  msg.MessageContent,
				"time":     time.Now().Format(time.RFC3339),
			})

		// ===============================
		//  HANDLE FOLLOW NOTIFICATION
		// ===============================
		msg.MessageContent="sent a message to the group"
		BrodcastGroupMembersNotification(msg.GroupID, currentUserID, map[string]any{
				"type":     "notification",
				"subType":  "group_message",
				"from":     currentUserID,
				"groupID":  msg.GroupID,
				"content":  msg.MessageContent,
				"name":     first_name + " " + last_name,
				"time":     time.Now().Format(time.RFC3339),
			})

			// ✅ Insert message direct sans check dyal follows
		case "follow":
			var followID int
			var first_name, last_name, photo, pryvsi string

			err = repository.Db.QueryRow(`SELECT first_name , last_name, image FROM users WHERE id = ?`, currentUserID).Scan(&first_name, &last_name, &photo)
			_ = repository.Db.QueryRow(`SELECT privacy FROM users WHERE id = ?`, msg.ReceiverId).Scan(&pryvsi)
			if err != nil {
				log.Println("DB error getting user info:", err)
				continue
			}

			query := `SELECT id FROM followers WHERE user_id = ? AND follower_id = ?`
			_ = repository.Db.QueryRow(query, msg.ReceiverId, currentUserID).Scan(&followID)
			if followID != 0 {
				msg.SubType = "unfollow"
				msg.First_name = first_name
				msg.Last_name = last_name

				msg.Photo = photo
				msg.MessageContent = "has unfollowed you"
			} else if followID == 0 && pryvsi == "public" {
				msg.SubType = "follow"
				msg.First_name = first_name
				msg.Last_name = last_name
				msg.Photo = photo
				msg.MessageContent = "has following you"
				q := `INSERT INTO notifications ( sender_id, receiver_id, type, message, created_at) VALUES (?, ?, ?, ?, ?) `
				_, _ = repository.Db.Exec(q, currentUserID, msg.ReceiverId, msg.Type, msg.MessageContent, time.Now().Unix())
			} else {
				continue
			}

			// Notify all connected users (except current user)
			BrodcastNotification(msg.ReceiverId, map[string]any{
				"type":       "notification",
				"subType":    msg.SubType,
				"from":       currentUserID,
				"first_name": msg.First_name,
				"last_name":  msg.Last_name,

				"photo":   msg.Photo,
				"content": msg.MessageContent,
				"time":    time.Now().Format(time.RFC3339),
			})
		case "invite_to_group":
			var first_name, last_name, photo string

			err = repository.Db.QueryRow(`SELECT first_name , last_name, image FROM users WHERE id = ?`, currentUserID).Scan(&first_name, &last_name, &photo)
			if err != nil {
				log.Println("DB error getting user info:", err)
				continue
			}
			msg.MessageContent="has invited you to join a group"
			q := `INSERT INTO notifications ( sender_id, receiver_id, type, message, created_at) VALUES (?, ?, ?, ?, ?) `
			_, _ = repository.Db.Exec(q, currentUserID, msg.ReceiverId, msg.Type, msg.MessageContent, time.Now().Unix())
			// Notify the invited user
			BrodcastNotification(msg.ReceiverId, map[string]any{
				"type":       "notification",
				"subType":    "group_invite",
				"from":       currentUserID,
				"first_name": first_name,
				"last_name":  last_name,
				"photo":      photo,
				"content":    msg.MessageContent,
				"time":       time.Now().Format(time.RFC3339),
			})

		default:
			log.Println("Unknown message type:", msg.Type)
		}

	}
}

func BrodcastGroupMembersNotification(groupID string, senderID string, message map[string]any) {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	rows, err := repository.Db.Query("SELEC user_id FROM group_members WHERE group_id = ?", groupID)
	if err != nil {
		log.Println("DB error getting group members:", err)
		return
	}
	defer rows.Close()

	var userID string
	for rows.Next() {
		if err := rows.Scan(&userID); err != nil {
			log.Println("DB error scanning group member:", err)
			continue
		}
		if userID == senderID {
			continue // Skip sender
		}
		conns, exists := Clients[userID]
		if !exists {
			continue
		}
		for _, conn := range conns {
			if err := conn.WriteJSON(message); err != nil {
				log.Println("WebSocket write error:", err)
				if err := conn.Close(); err != nil {
					log.Println("WebSocket close error:", err)
				}
			}
		}
	}
}

func BrodcastNotification(userID string, message map[string]any) {
	// TO BE IMPLEMENTED
	sendToUser(userID, message)
}

// BrodcastOnlineStatus notifies all connected clients about a user's online status change
func BrodcastOnlineStatus(userID string, online bool) {
	message := map[string]any{
		"type":   "logout",
		"userID": userID,
		"online": online,
	}
	for id, conns := range Clients {
		if id == userID {
			continue
		}
		for _, conn := range conns {
			if err := conn.WriteJSON(message); err != nil {
				log.Println("WebSocket write error:", err)
				if err := conn.Close(); err != nil {
					log.Println("WebSocket close error:", err)
				}
			}
		}
	}
}

// sendToUser sends a message to all WebSocket connections of a specific user
func sendToUser(userID string, message map[string]any) {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	conns, exists := Clients[userID]
	if !exists {
		return
	}

	for _, conn := range conns {
		if err := conn.WriteJSON(message); err != nil {
			log.Println("WebSocket write error:", err)
			if err := conn.Close(); err != nil {
				log.Println("WebSocket close error:", err)
			}
		}
	}
}

func sendToGroupMembers(groupID string, senderID string, message map[string]any)  {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	rows, err := repository.Db.Query("SELEC user_id FROM group_members WHERE group_id = ?", groupID)
	if err != nil {
		log.Println("DB error getting group members:", err)
		return
	}
	defer rows.Close()

	var userID string
	for rows.Next() {
		if err := rows.Scan(&userID); err != nil {
			log.Println("DB error scanning group member:", err)
			continue
		}
		if userID == senderID {
			continue // Skip sender
		}
		conns, exists := Clients[userID]
		if !exists {
			continue
		}
		for _, conn := range conns {
			if err := conn.WriteJSON(message); err != nil {
				log.Println("WebSocket write error:", err)
				if err := conn.Close(); err != nil {
					log.Println("WebSocket close error:", err)
				}
			}
		}
	}
}
