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
	Name           string `json:"name"`
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

		var nickname string
		err := repository.Db.QueryRow(`SELECT nickname FROM users WHERE id = ?`, currentUserID).Scan(&nickname)
		if err != nil {
			log.Println("DB error getting sender nickname:", err)
			nickname = "Unknown"
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
				"name":    nickname,
				"time":    time.Now().Format(time.RFC3339),
			})

		// ===============================
		//  HANDLE FOLLOW NOTIFICATION
		// ===============================
		case "follow":
			var followID int
			var name, photo string

			err := repository.Db.QueryRow(`SELECT nickname, image FROM users WHERE id = ?`, currentUserID).Scan(&name, &photo)
			if err != nil {
				log.Println("DB error getting user info:", err)
				continue
			}

			query := `SELECT id FROM followers WHERE user_id = ? AND follower_id = ?`
			_ = repository.Db.QueryRow(query, msg.ReceiverId, currentUserID).Scan(&followID)
			if followID != 0 {
				msg.SubType = "unfollow"
				msg.Name = name
				msg.Photo = photo
				msg.MessageContent = "has unfollowed you"
			} else {
				msg.SubType = "follow"
				msg.Name = name
				msg.Photo = photo
				msg.MessageContent = "has following you"
				q := `INSERT INTO notifications ( sender_id, receiver_id, type, message, created_at) VALUES (?, ?, ?, ?, ?) `
				_, _ = repository.Db.Exec(q, currentUserID, msg.ReceiverId, msg.Type, msg.MessageContent, time.Now().Unix())
			}

			// Notify all connected users (except current user)
			BrodcastNotification(msg.ReceiverId, map[string]any{
				"type":    "notification",
				"subType": msg.SubType,
				"from":    currentUserID,
				"name":    msg.Name,
				"photo":   msg.Photo,
				"content": msg.MessageContent,
				"time":    time.Now().Format(time.RFC3339),
			})
		default:
			log.Println("Unknown message type:", msg.Type)
		}

		log.Printf("Received message from user %s: %+v", currentUserID, msg)
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
