package handlers

import (
	"log"
	"net/http"
	"time"

	"social-network/internal/api/service"
	"social-network/internal/helper"
	"social-network/internal/repository/model"

	"github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections and manages user sessions
func Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := service.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	currentUserID, err := helper.AuthenticateUser(r)
	if err != nil {
		log.Println("Authentication error:", err)
		return
	}

	service.ClientsMutex.Lock()
	service.Clients[currentUserID] = append(service.Clients[currentUserID], conn)
	service.ClientsMutex.Unlock()

	service.BrodcastOnlineListe()
	// Listen for incoming messages
	Loop(conn, currentUserID)

	defer func() {
		service.ClientsMutex.Lock()
		conns, ok := service.Clients[currentUserID]

		if ok {
			for i, c := range conns {
				if c == conn {
					service.Clients[currentUserID] = append(conns[:i], conns[i+1:]...)
					break
				}
			}
			if len(service.Clients[currentUserID]) == 0 {
				delete(service.Clients, currentUserID)
				service.BrodcastOnlineStatus(currentUserID, false)
				if err := conn.Close(); err != nil {
					log.Println("WebSocket close error:", err)
				}
			}
		}
		service.ClientsMutex.Unlock()
	}()
}

func Loop(conn *websocket.Conn, currentUserID string) {
	for {
		var msg model.Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		user, err := model.GetUserByID(currentUserID)
		if err != nil {
			log.Println("DB error getting user info:", err)
			continue
		}

		msg.First_name = user["first_name"].(string)
		msg.Last_name = user["last_name"].(string)
		msg.Photo = user["photo"].(string)
		privaci := user["privacy"].(string)
		switch msg.Type {
		case "logout":
			service.BrodcastOnlineStatus(currentUserID, false)
		//  HANDLE CHAT MESSAGE
		case "online_list":
			// ✅ Send current online list to the new user
			service.BrodcastOnlineListe()
		case "message":

			if msg.ReceiverId == "" {
				log.Println("Invalid recipient ID")
				continue
			}
			valid, err := model.CheckIfUsersFollowEachOther(currentUserID, msg)
			if err != nil {
				log.Println("DB error checking follow status:", err)
				continue
			}
			if !valid {
				log.Println("Users do not follow each other")
				continue
			}
			// save notification
			err = model.SaveNotification(currentUserID, msg)
			if err != nil {
				log.Println("DB error saving notification:", err)
			}

			// save message
			err = model.SaveMessage(currentUserID, msg)
			if err != nil {
				log.Println("DB error saving message:", err)
				continue
			}

			// ✅ send message to receiver
			service.SendToUser(msg.ReceiverId, map[string]any{
				"type":    msg.Type,
				"from":    currentUserID,
				"to":      msg.ReceiverId,
				"content": msg.MessageContent,
				"time":    time.Now().Format(time.RFC3339),
			})

			// ✅ send back to sender also
			service.SendToUser(currentUserID, map[string]any{
				"type":    msg.Type,
				"from":    currentUserID,
				"to":      msg.ReceiverId,
				"content": msg.MessageContent,
				"time":    time.Now().Format(time.RFC3339),
			})

			// send notification to receiver
			service.BrodcastNotification(msg.ReceiverId, map[string]any{
				"type":    "notification",
				"subType": "message",
				"from":    currentUserID,
				"content": "sent you a message",
				"name":    msg.First_name + " " + msg.Last_name,
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

			valid, err := model.IsUserGroupMember(currentUserID, msg)
			if err != nil {
				log.Println("DB error checking group membership:", err)
				continue
			}
			if !valid {
				log.Println("User is not a member of the group")
				continue
			}

			// save message to DB
			err = model.SaveGroupMessage(currentUserID, msg)
			if err != nil {
				log.Println("DB error saving group message:", err)
				continue
			}

			service.SendToGroupMembers(msg.GroupID, currentUserID, map[string]any{
				"type":     msg.Type,
				"senderId": currentUserID,
				"groupID":  msg.GroupID,
				"content":  msg.MessageContent,
				"time":     time.Now().Format(time.RFC3339),
			})

			// ===============================
			//  HANDLE FOLLOW NOTIFICATION
			// ===============================
			msg.MessageContent = "sent a message to the group"
			service.BrodcastGroupMembersNotification(msg.GroupID, currentUserID, map[string]any{
				"type":    "notification",
				"subType": "group_message",
				"from":    currentUserID,
				"groupID": msg.GroupID,
				"content": msg.MessageContent,
				"name":    msg.First_name + " " + msg.Last_name,
				"time":    time.Now().Format(time.RFC3339),
			})

			// ✅ Insert message direct sans check dyal follows
		case "follow":

			exict, err := model.IsFollowingReceiver(currentUserID, msg)
			if err != nil {
				log.Println("DB error checking follow status:", err)
				continue
			}

			if exict {
				msg.SubType = "unfollow"
				msg.MessageContent = "has unfollowed you"
			} else if !exict && privaci == "public" {
				msg.SubType = "follow"
				msg.MessageContent = "has following you"
				err = model.SaveFollowNotification(currentUserID, msg)
				if err != nil {
					log.Println("DB error saving follow notification:", err)
					continue
				}
			} else {
				continue
			}

			// Notify all connected users (except current user)
			service.BrodcastNotification(msg.ReceiverId, map[string]any{
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

			msg.MessageContent = "has invited you to join a group"
			err := model.SaveGroupInvitationNotification(currentUserID, msg)
			if err != nil {
				log.Println("DB error saving group invitation notification:", err)
				continue
			}
			// Notify the invited user
			service.BrodcastNotification(msg.ReceiverId, map[string]any{
				"type":       "notification",
				"subType":    "group_invite",
				"from":       currentUserID,
				"first_name": msg.First_name,
				"last_name":  msg.Last_name,
				"photo":      msg.Photo,
				"content":    msg.MessageContent,
				"time":       time.Now().Format(time.RFC3339),
			})

		default:
			log.Println("Unknown message type:", msg.Type)
		}

	}
}
