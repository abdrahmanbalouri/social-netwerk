package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"social-network/app/api/service"
	"social-network/app/helper"
	"social-network/app/repository/model"

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
		fmt.Println("1cdd2f", msg.Type)
		switch msg.Type {
		case "logout":
			service.BrodcastOnlineStatus(currentUserID, false)
			return
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
			var imageFileName string
			// Si le message contient une image
			if msg.PictureSend != "" {
				imageFileName, err = helper.Upload(msg.PictureSend)
				if err != nil {
					log.Println("failed to upload image/file")
				}
			}
			// save message
			err = model.SaveMessage(currentUserID, msg, imageFileName)
			if err != nil {
				log.Println("DB error saving message:", err)
				continue
			}

			// ✅ send message to receiver
			service.SendToUser(msg.ReceiverId, map[string]any{
				"type":        msg.Type,
				"from":        currentUserID,
				"to":          msg.ReceiverId,
				"content":     msg.MessageContent,
				"time":        time.Now().Format(time.RFC3339),
				"name":        msg.First_name + " " + msg.Last_name,
				"image":       msg.Photo,
				"PictureSend": imageFileName,
			})

			// ✅ send back to sender also
			service.SendToUser(currentUserID, map[string]any{
				"type":        msg.Type,
				"from":        currentUserID,
				"to":          msg.ReceiverId,
				"content":     msg.MessageContent,
				"time":        time.Now().Format(time.RFC3339),
				"name":        msg.First_name + " " + msg.Last_name,
				"image":       msg.Photo,
				"PictureSend": imageFileName,
			})

			// send notification to receiver
			service.BrodcastNotification(msg.ReceiverId, map[string]any{
				"type":       "notification",
				"subType":    "message",
				"from":       currentUserID,
				"content":    "sent you a message",
				"first_name": msg.First_name,
				"last_name":  msg.Last_name,
				"time":       time.Now().Format(time.RFC3339),
				"image":      msg.Photo,
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
			var imageFileName string
			// Si le message contient une image
			if msg.PictureSend != "" {
				imageFileName, err = helper.Upload(msg.PictureSend)
				if err != nil {
					log.Println("failed to upload image/file")
				}
			}
			// save message to DB
			err = model.SaveGroupMessage(currentUserID, msg, imageFileName)
			if err != nil {
				log.Println("DB error saving group message:", err)
				continue
			}

			err = model.SaveGroupMessageNotification(currentUserID, msg)
			if err != nil {
				log.Println("DB error saving group message notification:", err)
			}

			service.SendToGroupMembers(msg.GroupID, currentUserID, map[string]any{
				"type":        msg.Type,
				"from":        currentUserID,
				"groupID":     msg.GroupID,
				"content":     msg.MessageContent,
				"time":        time.Now().Format(time.RFC3339),
				"first_name":  msg.First_name,
				"last_name":   msg.Last_name,
				"image":       msg.Photo,
				"PictureSend": imageFileName,
			})

			// ===============================
			//  HANDLE FOLLOW NOTIFICATION
			// ===============================
			msg.MessageContent = "sent a message to the group"
			service.BrodcastGroupMembersNotification(msg.GroupID, currentUserID, map[string]any{
				"type":       "notification",
				"subType":    "group_message",
				"from":       currentUserID,
				"groupID":    msg.GroupID,
				"content":    msg.MessageContent,
				"first_name": msg.First_name,
				"last_name":  msg.Last_name,
				"time":       time.Now().Format(time.RFC3339),
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
				"photo":      msg.Photo,
				"content":    msg.MessageContent,
				"time":       time.Now().Format(time.RFC3339),
			})
		case "followRequest":
			exict, err := model.IsFollowingReceiver(currentUserID, msg)
			if err != nil {
				log.Println("DB error checking follow status:", err)
				continue
			}

			if !exict {
				msg.SubType = "followRequest"
				msg.MessageContent = "send you a followRequest"
				err = model.SaveFollowNotification(currentUserID, msg)
				if err != nil {
					log.Println("DB error saving follow notification:", err)
					continue
				}
			} else {
				msg.SubType = "unfollow"
				msg.MessageContent = "has unfollowed you"
			}

			// Notify all connected users (except current user)
			service.BrodcastNotification(msg.ReceiverId, map[string]any{
				"type":       "notification",
				"subType":    msg.SubType,
				"from":       currentUserID,
				"first_name": msg.First_name,
				"last_name":  msg.Last_name,
				"photo":      msg.Photo,
				"content":    msg.MessageContent,
				"time":       time.Now().Format(time.RFC3339),
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
		case "joinRequest":

			msg.MessageContent = "has requested to join your group"
			err, receiver := model.SaveGroupJoinRequestNotification(currentUserID, msg)
			if err != nil {
				log.Println("DB error saving group join request notification:", err)
				continue
			}
			// Notify the group admin
			service.BrodcastNotification(receiver, map[string]any{
				"type":       "notification",
				"subType":    "group_join_request",
				"from":       currentUserID,
				"first_name": msg.First_name,
				"last_name":  msg.Last_name,
				"photo":      msg.Photo,
				"content":    msg.MessageContent,
				"time":       time.Now().Format(time.RFC3339),
			})
		case "new event":
			fmt.Println("1111111111111")
			msg.MessageContent = "has a new event"
			err := model.SaveGroupInvitationNotification(currentUserID, msg)
			if err != nil {
				log.Println("DB error saving group invitation notification:", err)
				continue
			}
			err = model.SaveGroupMessageNotification(currentUserID, msg)
			if err != nil {
				log.Println("DB error saving group message notification:", err)
			}
			groupname, err := model.Name(msg)
			if err != nil {
				log.Println("DB error saving group invitation notification:", err)
				continue
			}
			fmt.Println("11", msg.ReceiverId, "22", currentUserID)
			service.BrodcastGroupMembersNotification(msg.ReceiverId, currentUserID, map[string]any{
				"type":       "notification",
				"subType":    "group_message",
				"from":       currentUserID,
				"groupID":    msg.GroupID,
				"content":    msg.MessageContent,
				"first_name": groupname,
				"last_name":  "",
				"time":       time.Now().Format(time.RFC3339),
			})
		default:
			log.Println("Unknown message type:", msg.Type)
		}

	}
}
