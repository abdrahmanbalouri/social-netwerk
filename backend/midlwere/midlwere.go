package middlewares

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

func SessionMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedPaths := []string{
			"/api/login",
			"/api/register",
			"/api/logout",
		}
		
		for _, path := range allowedPaths {
			if r.URL.Path == path {
				fmt.Println(r.URL.Path)
				next.ServeHTTP(w, r)
				return
			}
		}

		fmt.Println(r.URL.Path)


		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			http.Error(w, "unauthorized: no session cookie", http.StatusUnauthorized)
			return
		}

		var userID string
		var expiresAt time.Time
		err = db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE token = ?", cookie.Value).
			Scan(&userID, &expiresAt)
		if err != nil {
			http.Error(w, "unauthorized: invalid session", http.StatusUnauthorized)
			return
		}

		if time.Now().After(expiresAt) {
			http.Error(w, "session expired", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
