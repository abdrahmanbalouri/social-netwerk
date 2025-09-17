package handlers

import (
	"fmt"
	"net/http"

	"social-network/internal/repository"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"unauthorized"}`))
		return

	}
	_, err = repository.Db.Query("SELECT user_id FROM sessions WHERE token=?", c.Value)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"unauthorized"}`))
		return
	} else {
		fmt.Println("authorized")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"authorized"}`))
		return
	}
}
