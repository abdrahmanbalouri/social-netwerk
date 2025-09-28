package api

import (
	"database/sql"
	"net/http"

	handlers "social-network/internal/api/handler"
)

func Routes(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
	})
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	mux.HandleFunc("/api/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/login", handlers.LoginHandler)
	mux.HandleFunc("/api/logout", handlers.LogoutHandler)
	mux.HandleFunc("/api/me", handlers.MeHandler)
	mux.HandleFunc("/api/profile", handlers.ProfileHandler)
	mux.HandleFunc("/api/createpost", handlers.Createpost)
	mux.HandleFunc("/api/Getpost/{id}", handlers.GetPostsHandler)
	mux.HandleFunc("/api/Getallpost", handlers.AllpostsHandler)
	mux.HandleFunc("/api/GetUsersHandler", handlers.GetUsersHandler)
	mux.HandleFunc("/api/Getcomments/{id}", handlers.GetCommentsHandler)
	mux.HandleFunc("/api/createcomment", handlers.CreateCommentHandler)
	mux.HandleFunc("/api/groups", handlers.AddGroupHandler)
	mux.HandleFunc("/invitations/respond", handlers.GroupInvitationResponse)
	mux.HandleFunc("/group/invitation", handlers.GroupInvitationRequest)
	mux.HandleFunc("/group/addPost", handlers.CreatePostGroup)

	return mux
}
