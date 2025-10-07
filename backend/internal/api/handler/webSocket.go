package handlers

import (
	"log"
	"net/http"
	"sync"

	"social-network/internal/helper"

	"github.com/gorilla/websocket"
)

// Global WebSocket-related variables managed within this package
var (
	Clients      = make(map[string][]*websocket.Conn) // Use models.Client
	ClientsMutex sync.Mutex
	Upgrader     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
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
		case "message":
			// Handle chat message broadcasting
			recipientID, ok := msg["to"].(string)
			if !ok {
				log.Println("Invalid recipient ID")
				continue
			}
			sendToUser(recipientID, msg)
			sendToUser(currentUserID, msg) // Echo back to sender
		default:
			log.Println("Unknown message type:", typeMsg)
		}
		log.Printf("Received message from user %s: %v", currentUserID, msg)
	}
	defer func() {
		ClientsMutex.Lock()
		conns := Clients[currentUserID]
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
