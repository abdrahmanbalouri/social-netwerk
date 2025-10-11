package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"sync"

	"social-network/internal/helper"
	"social-network/internal/repository"

	"github.com/gorilla/websocket"
)

// Global WebSocket-related variables managed within this package
var (
	Clients      = make(map[string][]*websocket.Conn) // Use models.Client
	ClientsMutex sync.Mutex
	Upgrader     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow connections from localhost:3000
			return r.Header.Get("Origin") == "http://localhost:3000"
		},
	}
)

// WebSocketHandler handles WebSocket connections and manages user sessions
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
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
		var msg map[string]any
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
		typeMsg, ok := msg["type"].(string)
		if !ok {
			log.Println("Invalid message format")
			continue
		}
		switch typeMsg {
		case "logout":
			// Handle user logout
			BrodcastOnlineStatus(currentUserID, false)
			delete(Clients, currentUserID)
			return // Exit the loop to close the connection
		case "message":
			// Handle chat message broadcasting
			recipientID, ok := msg["to"].(string)
			if !ok {
				log.Println("Invalid recipient ID")
				continue
			}
			// Check if recipient exists and if either user follows the other
			var exists int
			err := repository.Db.QueryRow(`
				SELECT 1 FROM users WHERE id = ?
			`, recipientID).Scan(&exists)
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
			`, currentUserID, recipientID, recipientID, currentUserID).Scan(&followExists)
			if err == sql.ErrNoRows {
				log.Println("No follow relationship between users")
				continue
			} else if err != nil {
				log.Println("DB error:", err)
				continue
			}
			// Store message in the database and send to recipient
			_, err = repository.Db.Exec("INSERT INTO messages (sender_id, recipient_id, content) VALUES (?, ?, ?)", currentUserID, recipientID, msg["content"])
			if err != nil {
				log.Println("DB error:", err)
				continue
			}
			// Send the message to the recipient and the sender
			sendToUser(recipientID, msg)
			sendToUser(currentUserID, msg) // Echo back to sender
		default:
			log.Println("Unknown message type:", typeMsg)
		}
		log.Printf("Received message from user %s: %v", currentUserID, msg)
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
