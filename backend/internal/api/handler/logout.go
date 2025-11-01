package handlers

import (
	"net/http"
	"time"

	"social-network/internal/repository"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err == nil {
		_, dbErr := repository.Db.Exec("DELETE FROM sessions WHERE token=?", c.Value)
		if dbErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Database error during logout"))
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Unix(0, 0),
		})
		w.Write([]byte("logged out"))
	}
}
