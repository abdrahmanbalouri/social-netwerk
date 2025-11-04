package handlers

import (
	"net/http"
	"time"

	service "social-network/internal/api/sevice"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("no session"))
		return
	}

	err = service.LogoutUser(c.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("logout failed"))
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
