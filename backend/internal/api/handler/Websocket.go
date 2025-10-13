package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"

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
	ReceiverId     string `json:"receiverId"`
	MessageContent string `json:"messageContent"`
	Name           string `json:"name"`
	Photo          string `json:"photo"`
	Content        string `json:"content"`
	To             string `json:"to"`
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
	BrodcastOnlineStatus(currentUserID, true)

	// Listen for incoming messages
	Loop(conn, currentUserID)

	defer func() {
		ClientsMutex.Lock()
		conns, ok := Clients[currentUserID]
		if !ok {
			ClientsMutex.Unlock()
			return
		}
		for i, c := range conns {
			if c == conn {
				Clients[currentUserID] = append(conns[:i], conns[i+1:]...)
				break
			}
		}
		if len(Clients[currentUserID]) == 0 {
			delete(Clients, currentUserID)
			BrodcastOnlineStatus(currentUserID, false)
		}
		ClientsMutex.Unlock()
		conn.Close()
		log.Printf("User %s disconnected", currentUserID)
	}()
}

func Loop(conn *websocket.Conn, currentUserID string) {
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		switch msg.Type {
		// ===============================
		//  HANDLE LOGOUT
		// ===============================
		case "logout":
			BrodcastOnlineStatus(currentUserID, false)
			delete(Clients, currentUserID)
			return

		// ===============================
		//  HANDLE CHAT MESSAGE
		// ===============================
		case "message":
			if msg.To == "" {
				log.Println("Invalid recipient ID")
				continue
			}

			var exists int
			err := repository.Db.QueryRow(`
				SELECT 1 FROM users WHERE id = ?
			`, msg.To).Scan(&exists)
			if err == sql.ErrNoRows {
				log.Println("Recipient user not found")
				continue
			} else if err != nil {
				log.Println("DB error:", err)
				continue
			}

			var followExists int
			err = repository.Db.QueryRow(`
				SELECT 1 FROM follows 
				WHERE (user_id = ? AND followed_id = ?)
				   OR (user_id = ? AND followed_id = ?)
			`, currentUserID, msg.To, msg.To, currentUserID).Scan(&followExists)
			if err == sql.ErrNoRows {
				log.Println("No follow relationship between users")
				continue
			} else if err != nil {
				log.Println("DB error:", err)
				continue
			}

			_, err = repository.Db.Exec("INSERT INTO messages (sender_id, recipient_id, content) VALUES (?, ?, ?)", currentUserID, msg.To, msg.Content)
			if err != nil {
				log.Println("DB error:", err)
				continue
			}

			sendToUser(msg.To, map[string]any{
				"type":    "message",
				"from":    currentUserID,
				"content": msg.Content,
			})
			sendToUser(currentUserID, map[string]any{
				"type":    "message",
				"from":    currentUserID,
				"content": msg.Content,
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
			_ = repository.Db.QueryRow(query, currentUserID, msg.ReceiverId).Scan(&followID)

			if followID != 0 {
				msg.Type = "already_following"
			} else {
				msg.Name = name
				msg.Photo = photo
				msg.MessageContent = "has followed you"
			}

			// Notify all connected users (except current user)
			ClientsMutex.Lock()
			for uid, conns := range Clients {
				if uid == currentUserID {
					continue
				}
				for _, con := range conns {
					jsonMsg, err := json.Marshal(msg)
					if err != nil {
						continue
					}
					if err := con.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
						log.Println("WebSocket write error:", err)
					}
				}
			}
			ClientsMutex.Unlock()

		default:
			log.Println("Unknown message type:", msg.Type)
		}

		log.Printf("Received message from user %s: %+v", currentUserID, msg)
	}
}

// BrodcastOnlineStatus notifies all connected clients about a user's online status change
func BrodcastOnlineStatus(userID string, online bool) {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	message := map[string]any{
		"type":   "status",
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
				conn.Close()
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
			conn.Close()
		}
	}
}
