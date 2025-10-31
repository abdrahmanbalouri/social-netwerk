package api

import (
	"net/http"

	handlers "social-network/internal/api/handler"
)

func Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
	})
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	mux.HandleFunc("/api/followers/", handlers.FollowersHandler)
	mux.HandleFunc("/api/followRequest", handlers.FollowRequest)
	mux.HandleFunc("/api/groupeInvitation", handlers.GroupeInvitation)
	mux.HandleFunc("/api/getEvents/", handlers.GetEvents)
	mux.HandleFunc("/api/createEvent/", handlers.CreateEvent)
	mux.HandleFunc("/api/event/action/", handlers.EventAction)







	mux.HandleFunc("/api/followRequest/action", handlers.FollowRequestAction)
	mux.HandleFunc("/api/following/", handlers.FollowingHandler)
	mux.HandleFunc("/api/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/login", handlers.LoginHandler)
	mux.HandleFunc("/api/logout", handlers.LogoutHandler)
	mux.HandleFunc("/api/me", handlers.MeHandler)
	mux.HandleFunc("/api/profile", handlers.ProfileHandler)
	mux.HandleFunc("/api/createpost", handlers.Createpost)
	mux.HandleFunc("/api/Getpost/{id}", handlers.GetPostsHandler)
	mux.HandleFunc("/api/Getallpost/{id}", handlers.AllpostsHandler)
	mux.HandleFunc("/api/getmypost/{id}/{offset}", handlers.Getmypost)
	mux.HandleFunc("/api/GetUsersHandler", handlers.GetUsersHandler)
	mux.HandleFunc("/api/communfriends", handlers.GetCommunFriends)
	mux.HandleFunc("/api/Getcomments/{id}/{offset}", handlers.GetCommentsHandler)
	mux.HandleFunc("/api/gallery", handlers.GalleryHandler)
	mux.HandleFunc("/api/createcomment", handlers.CreateCommentHandler)
	mux.HandleFunc("/api/editor", handlers.Editor)
	mux.HandleFunc("/api/like/{id}", handlers.LikeHandler)
	mux.HandleFunc("/api/follow", handlers.FollowHandler)
	mux.HandleFunc("/api/users/followers", handlers.Getfollowers)
	mux.HandleFunc("/api/getlastcomment/{id}", handlers.Getlastcommnet)
	mux.HandleFunc("/api/searchUser", handlers.SearchUserHandler)
	mux.HandleFunc("/notifcation", handlers.Notifications)
	mux.HandleFunc("/api/Getstories", handlers.GetStories)
	mux.HandleFunc("/api/Createstories", handlers.CreateStories)
	mux.HandleFunc("/api/getmessages", handlers.GetMessagesHandler)
	mux.HandleFunc("/ws", handlers.Websocket)
	mux.HandleFunc("/api/clearNotifications", handlers.ClearNotifications)
	mux.HandleFunc("/myGroups", handlers.GetMyGroups)
	mux.HandleFunc("/groups", handlers.GetAllGroups)
	mux.HandleFunc("/group/like", handlers.LikesGroup)
	mux.HandleFunc("/group/fetchComments", handlers.GetCommentGroup)
	mux.HandleFunc("/api/groups/add", handlers.CreateGroupHandler)
	mux.HandleFunc("/invitations/respond", handlers.GroupInvitationResponse)
	mux.HandleFunc("/group/invitation/{id}", handlers.GroupInvitationRequest)
	mux.HandleFunc("/group/addPost/{id}", handlers.CreatePostGroup)
	mux.HandleFunc("/group/fetchPosts/{id}", handlers.GetAllPostsGroup)
	mux.HandleFunc("/group/fetchPost/{id}", handlers.GetPostGroup)
	mux.HandleFunc("/group/addComment", handlers.CreateCommentGroup)
	mux.HandleFunc("/api/videos", handlers.GetVedioHandler)

	return mux
}

// {
//     "title": "Group 1",
//     "description": "group for developpers",
//     "invitedUsers": []
// }
