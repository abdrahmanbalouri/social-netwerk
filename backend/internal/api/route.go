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

	mux.HandleFunc("/api/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/login", handlers.LoginHandler)
	mux.HandleFunc("/api/logout", handlers.LogoutHandler)
	mux.HandleFunc("/api/me", handlers.MeHandler)
	mux.HandleFunc("/api/createpost", handlers.Createpost)

	return mux
}
