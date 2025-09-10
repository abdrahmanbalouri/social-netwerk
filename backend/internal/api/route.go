package api

import (
	"database/sql"
	"net/http"
)

func Routes(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Social Network API 🚀"))
	})

	// mux.HandleFunc("/register", RegisterHandler(db))
	// mux.HandleFunc("/login", LoginHandler(db))
	// mux.HandleFunc("/logout", LogoutHandler(db))

	return mux
}
