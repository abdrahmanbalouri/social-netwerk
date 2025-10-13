package handlers

import (
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

	// Perform the query
	rows, err := repository.Db.Query("SELECT user_id FROM sessions WHERE token=?", c.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"unauthorized"}`))
		return
	}
	defer rows.Close() // Close the rows once we're done

	if !rows.Next() { // Check if we have at least one row
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"unauthorized"}`))
		return
	}

	// var userId int
	// if err := rows.Scan(&userId); err != nil {
	//     w.WriteHeader(http.StatusUnauthorized)
	//     w.Write([]byte(`{"message":"unauthorized"}`))
	//     return
	// }

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"authorized"}`))
}
