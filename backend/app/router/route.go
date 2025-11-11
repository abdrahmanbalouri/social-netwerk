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
	mux.HandleFunc("/api/register", middlewares.RateLimitLoginMiddleware(handlers.RegisterHandler))
	mux.HandleFunc("/api/login",  middlewares.RateLimitLoginMiddleware(handlers.LoginHandler))     
	mux.HandleFunc("/api/logout", handlers.LogoutHandler)    
	mux.HandleFunc("/api/me", handlers.MeHandler)            
	mux.HandleFunc("/api/profile", handlers.ProfileHandler)
	mux.HandleFunc("/api/searchUser", handlers.SearchUserHandler)   
	mux.HandleFunc("/api/GetUsersHandler", handlers.GetUsersHandler)

	// ======= Followers / Following =======
	mux.HandleFunc("/api/followers/", handlers.FollowersHandler)             
	mux.HandleFunc("/api/following/", handlers.FollowingHandler)             
	mux.HandleFunc("/api/followRequest", handlers.FollowRequest)             
	mux.HandleFunc("/api/followRequest/action", handlers.FollowRequestAction)
	mux.HandleFunc("/api/follow",  middlewares.RatelimitMiddleware(handlers.FollowHandler, "follow", 30))
	mux.HandleFunc("/api/users/followers", handlers.Getfollowers)            
	mux.HandleFunc("/api/communfriends", handlers.GetCommunFriends)

	// ======= Posts & Comments =======
	mux.HandleFunc("/api/createpost", middlewares.RatelimitMiddleware(handlers.Createpost, "posts", 30))         
	mux.HandleFunc("/api/Getpost/{id}", handlers.GetPostsHandler)                
	mux.HandleFunc("/api/Getallpost/{id}", handlers.AllpostsHandler)             
	mux.HandleFunc("/api/getmypost/{id}/{offset}", handlers.Getmypost)           
	mux.HandleFunc("/api/createcomment", middlewares.RatelimitMiddleware(handlers.CreateCommentHandler, "comments", 80))      
	mux.HandleFunc("/api/Getcomments/{id}/{offset}", handlers.GetCommentsHandler)
	mux.HandleFunc("/api/like/{id}", middlewares.RatelimitMiddleware(handlers.LikeHandler, "likes", 40))              
	mux.HandleFunc("/api/getlastcomment/{id}", handlers.Getlastcommnet)          
	mux.HandleFunc("/api/editor", handlers.Editor)                               
	mux.HandleFunc("/api/gallery", handlers.GalleryHandler)                      

	// ======= Stories & Media =======
	mux.HandleFunc("/api/Getstories", handlers.GetStories)      
	mux.HandleFunc("/api/Createstories", middlewares.RatelimitMiddleware(handlers.CreateStories, "story",40))
	mux.HandleFunc("/api/videos", handlers.GetVideoHandler)     

	// ======= Notifications =======
	mux.HandleFunc("/notifcation", handlers.Notifications)
	mux.HandleFunc("/api/clearNotifications", handlers.ClearNotifications)

	// ======= Messaging / WebSocket =======
	mux.HandleFunc("/api/getmessages", handlers.GetMessagesHandler)          
	mux.HandleFunc("/api/getGroupMessages", handlers.GetGroupMessagesHandler)
	mux.HandleFunc("/ws", handlers.Websocket)                                

	// ======= Groups =======
	mux.HandleFunc("/myGroups", handlers.GetMyGroups)                        
	mux.HandleFunc("/groups", handlers.GetAllGroups)                         
	mux.HandleFunc("/api/groups/add", middlewares.RatelimitMiddleware(handlers.CreateGroupHandler, "createGroup", 50))           
	mux.HandleFunc("/invitations/respond", handlers.GroupInvitationResponse) 
	mux.HandleFunc("/group/invitation/{id}", middlewares.RatelimitMiddleware(handlers.GroupInvitationRequest, "groupInvite", 100))
	mux.HandleFunc("/group/addPost/{id}", middlewares.RatelimitMiddleware(handlers.CreatePostGroupHandler, "posts", 30))   
	mux.HandleFunc("/group/fetchPosts/{id}", handlers.GetAllPostsGroup)      
	// mux.HandleFunc("/group/fetchPost/{id}", handlers.GetPostGroup)
	mux.HandleFunc("/group/like/{id}/{groupId}", middlewares.RatelimitMiddleware(handlers.LikesGroup, "likes", 40))                      
	mux.HandleFunc("/group/updatepost/{id}/{groupId}", handlers.GetGroupPostByID)          
	mux.HandleFunc("/group/Getcomments/{id}/{offset}/{groupId}", handlers.GetCommentsGroup)
	mux.HandleFunc("/group/getlastcomment/{id}/{groupId}", handlers.GetlastcommnetGroup)   
	mux.HandleFunc("/api/fetchJoinRequests/{id}", handlers.FetchJoinRequests)              
	mux.HandleFunc("/api/fetchGroupInvitation", handlers.GroupeInvitation)                 
	mux.HandleFunc("/api/fetchFriendsForGroups/{id}", handlers.FetchFriendsForGroups)

	// ======= Events =======
	mux.HandleFunc("/api/getEvents/", handlers.GetEvents)     
	mux.HandleFunc("/api/createEvent/", middlewares.RatelimitMiddleware(handlers.CreateEvent, "events", 30))
	mux.HandleFunc("/api/event/action/", middlewares.RatelimitMiddleware(handlers.EventAction, "eventAction", 40))
	mux.HandleFunc("/api/myevents", handlers.MyEavents)

	return mux
}
