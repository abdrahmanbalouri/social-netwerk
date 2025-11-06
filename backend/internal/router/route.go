package api

import (
	"net/http"

	handlers "social-network/internal/api/handler"
)

func Routes() http.Handler {
	mux := http.NewServeMux()

	// ======= Home / Default =======
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
	})

	// ======= Static Files =======
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads")))) // hada mamahaloho mina i3rab

	// ======= Auth & User =======
	mux.HandleFunc("/api/register", handlers.RegisterHandler) // dart
	mux.HandleFunc("/api/login", handlers.LoginHandler)       // dart
	mux.HandleFunc("/api/logout", handlers.LogoutHandler)     // dart
	mux.HandleFunc("/api/me", handlers.MeHandler)             // dart
	mux.HandleFunc("/api/profile", handlers.ProfileHandler)
	mux.HandleFunc("/api/searchUser", handlers.SearchUserHandler)    // dart
	mux.HandleFunc("/api/GetUsersHandler", handlers.GetUsersHandler) // dart

	// ======= Followers / Following =======
	mux.HandleFunc("/api/followers/", handlers.FollowersHandler)              // dart
	mux.HandleFunc("/api/following/", handlers.FollowingHandler)              // dart
	mux.HandleFunc("/api/followRequest", handlers.FollowRequest)              // dart
	mux.HandleFunc("/api/followRequest/action", handlers.FollowRequestAction) // dart
	mux.HandleFunc("/api/follow", handlers.FollowHandler)                     //  dart
	mux.HandleFunc("/api/users/followers", handlers.Getfollowers)             // dart
	mux.HandleFunc("/api/communfriends", handlers.GetCommunFriends)

	// ======= Posts & Comments =======
	mux.HandleFunc("/api/createpost", handlers.Createpost)                        // dart
	mux.HandleFunc("/api/Getpost/{id}", handlers.GetPostsHandler)                 // dart
	mux.HandleFunc("/api/Getallpost/{id}", handlers.AllpostsHandler)              // dart
	mux.HandleFunc("/api/getmypost/{id}/{offset}", handlers.Getmypost)            // dart
	mux.HandleFunc("/api/createcomment", handlers.CreateCommentHandler)           // dart
	mux.HandleFunc("/api/Getcomments/{id}/{offset}", handlers.GetCommentsHandler) // dart
	mux.HandleFunc("/api/like/{id}", handlers.LikeHandler)                        // dart
	mux.HandleFunc("/api/getlastcomment/{id}", handlers.Getlastcommnet)           // dart
	mux.HandleFunc("/api/editor", handlers.Editor)                                // dart
	mux.HandleFunc("/api/gallery", handlers.GalleryHandler)                       // dart

	// ======= Stories & Media =======
	mux.HandleFunc("/api/Getstories", handlers.GetStories)       // dart
	mux.HandleFunc("/api/Createstories", handlers.CreateStories) // dart
	mux.HandleFunc("/api/videos", handlers.GetVideoHandler)      // dart

	// ======= Notifications =======
	mux.HandleFunc("/notifcation", handlers.Notifications)
	mux.HandleFunc("/api/clearNotifications", handlers.ClearNotifications)

	// ======= Messaging / WebSocket =======
	mux.HandleFunc("/api/getmessages", handlers.GetMessagesHandler) // dart
	mux.HandleFunc("/api/getGroupMessages", handlers.GetGroupMessagesHandler) // dart
	mux.HandleFunc("/ws", handlers.Websocket) // dart

	// ======= Groups =======
	mux.HandleFunc("/myGroups", handlers.GetMyGroups)
	mux.HandleFunc("/groups", handlers.GetAllGroups)
	mux.HandleFunc("/api/groups/add", handlers.CreateGroupHandler)
	mux.HandleFunc("/invitations/respond", handlers.GroupInvitationResponse)
	mux.HandleFunc("/group/invitation/{id}", handlers.GroupInvitationRequest)
	mux.HandleFunc("/group/addPost/{id}", handlers.CreatePostGroupHandler) // dart
	mux.HandleFunc("/group/fetchPosts/{id}", handlers.GetAllPostsGroup)    // dart
	// mux.HandleFunc("/group/fetchPost/{id}", handlers.GetPostGroup)
	mux.HandleFunc("/group/like/{id}/{groupId}", handlers.LikesGroup)                       // dart
	mux.HandleFunc("/group/updatepost/{id}/{groupId}", handlers.GetGroupPostByID)           // dart
	mux.HandleFunc("/group/Getcomments/{id}/{offset}/{groupId}", handlers.GetCommentsGroup) // dart
	mux.HandleFunc("/group/getlastcomment/{id}/{groupId}", handlers.GetlastcommnetGroup)    // dart
	mux.HandleFunc("/api/fetchJoinRequests/{id}", handlers.FetchJoinRequests)
	mux.HandleFunc("/api/fetchGroupInvitation", handlers.GroupeInvitation)
	mux.HandleFunc("/api/fetchFriendsForGroups/{id}", handlers.FetchFriendsForGroups)

	// ======= Events =======
	mux.HandleFunc("/api/getEvents/", handlers.GetEvents)      // dart
	mux.HandleFunc("/api/createEvent/", handlers.CreateEvent)  // dart
	mux.HandleFunc("/api/event/action/", handlers.EventAction) // dart
	mux.HandleFunc("/api/myevents", handlers.MyEavents)

	return mux
}
