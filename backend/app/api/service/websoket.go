package service

import (
	"log"
	"net/http"
	"sync"

	"social-network/app/repository/model"

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

func BrodcastGroupMembersNotification(groupID string, senderID string, message map[string]any) {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()
	
	groupMembers, err := model.GetGroupMembers(groupID)
	if err != nil {
		log.Println("DB error getting group members:", err)
		return
	}

	for _, userID := range groupMembers {
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
	SendToUser(userID, message)
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
func SendToUser(userID string, message map[string]any) {
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

func SendToGroupMembers(groupID string, senderID string, message map[string]any) {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	groupMembers, err := model.GetGroupMembers(groupID)
	if err != nil {
		log.Println("DB error getting group members:", err)
		return
	}

	for _, userID := range groupMembers {
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
