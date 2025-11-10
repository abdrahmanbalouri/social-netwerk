package api

import (
	"net/http"
	handlers "social-network/app/api/handler"
	middlewares "social-network/middleware"
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
	mux.HandleFunc("/api/follow",  middlewares.RatelimitMiddleware(handlers.FollowHandler, "follow", 3))               //  dart
	mux.HandleFunc("/api/users/followers", handlers.Getfollowers)             // dart
	mux.HandleFunc("/api/communfriends", handlers.GetCommunFriends)

	// ======= Posts & Comments =======
	mux.HandleFunc("/api/createpost", middlewares.RatelimitMiddleware(handlers.Createpost, "posts", 3))             //////////           // dart
	mux.HandleFunc("/api/Getpost/{id}", handlers.GetPostsHandler)                 // dart
	mux.HandleFunc("/api/Getallpost/{id}", handlers.AllpostsHandler)              // dart
	mux.HandleFunc("/api/getmypost/{id}/{offset}", handlers.Getmypost)            // dart
	mux.HandleFunc("/api/createcomment", middlewares.RatelimitMiddleware(handlers.CreateCommentHandler, "comments", 3))  //////        // dart
	mux.HandleFunc("/api/Getcomments/{id}/{offset}", handlers.GetCommentsHandler) // dart
	mux.HandleFunc("/api/like/{id}", middlewares.RatelimitMiddleware(handlers.LikeHandler, "likes", 3))         //////               // dart
	mux.HandleFunc("/api/getlastcomment/{id}", handlers.Getlastcommnet)           // dart
	mux.HandleFunc("/api/editor", handlers.Editor)                                // dart
	mux.HandleFunc("/api/gallery", handlers.GalleryHandler)                       // dart

	// ======= Stories & Media =======
	mux.HandleFunc("/api/Getstories", handlers.GetStories)       // dart
	mux.HandleFunc("/api/Createstories", middlewares.RatelimitMiddleware(handlers.CreateStories, "story", 3)) // dart ////////
	mux.HandleFunc("/api/videos", handlers.GetVideoHandler)      // dart

	// ======= Notifications =======
	mux.HandleFunc("/notifcation", handlers.Notifications)
	mux.HandleFunc("/api/clearNotifications", handlers.ClearNotifications)

	// ======= Messaging / WebSocket =======
	mux.HandleFunc("/api/getmessages", handlers.GetMessagesHandler)           // dart
	mux.HandleFunc("/api/getGroupMessages", handlers.GetGroupMessagesHandler) // dart
	mux.HandleFunc("/ws", handlers.Websocket)                                 // dart

	// ======= Groups =======
	mux.HandleFunc("/myGroups", handlers.GetMyGroups)                         // dart
	mux.HandleFunc("/groups", handlers.GetAllGroups)                          // dart
	mux.HandleFunc("/api/groups/add", middlewares.RatelimitMiddleware(handlers.CreateGroupHandler, "createGroup", 3))            // dart
	mux.HandleFunc("/invitations/respond", handlers.GroupInvitationResponse)  // dart
	mux.HandleFunc("/group/invitation/{id}", middlewares.RatelimitMiddleware(handlers.GroupInvitationRequest, "groupInvite", 3)) // dart
	mux.HandleFunc("/group/addPost/{id}", middlewares.RatelimitMiddleware(handlers.CreatePostGroupHandler, "posts", 3))    // dart
	mux.HandleFunc("/group/fetchPosts/{id}", handlers.GetAllPostsGroup)       // dart
	// mux.HandleFunc("/group/fetchPost/{id}", handlers.GetPostGroup)
	mux.HandleFunc("/group/like/{id}/{groupId}", middlewares.RatelimitMiddleware(handlers.LikesGroup, "likes", 2))                       // dart
	mux.HandleFunc("/group/updatepost/{id}/{groupId}", handlers.GetGroupPostByID)           // dart
	mux.HandleFunc("/group/Getcomments/{id}/{offset}/{groupId}", handlers.GetCommentsGroup) // dart
	mux.HandleFunc("/group/getlastcomment/{id}/{groupId}", handlers.GetlastcommnetGroup)    // dart
	mux.HandleFunc("/api/fetchJoinRequests/{id}", handlers.FetchJoinRequests)               // dart
	mux.HandleFunc("/api/fetchGroupInvitation", handlers.GroupeInvitation)                  // dart
	mux.HandleFunc("/api/fetchFriendsForGroups/{id}", handlers.FetchFriendsForGroups)

	// ======= Events =======
	mux.HandleFunc("/api/getEvents/", handlers.GetEvents)      // dart
	mux.HandleFunc("/api/createEvent/", middlewares.RatelimitMiddleware(handlers.CreateEvent, "events", 3)) ///////            // dart
	mux.HandleFunc("/api/event/action/", middlewares.RatelimitMiddleware(handlers.EventAction, "eventAction", 3)) // dart ///////////////
	mux.HandleFunc("/api/myevents", handlers.MyEavents)

	return mux
}
