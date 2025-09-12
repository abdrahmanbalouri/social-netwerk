package handlers

import (
	"net/http"
	"time"
	"social-network/internal/repository"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err == nil {
		repository.Db.Exec("DELETE FROM sessions WHERE token=?", c.Value)
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Unix(0, 0),
		})
	}
	w.Write([]byte("logged out"))
}
