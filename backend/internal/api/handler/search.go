// handlers/search.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

type UserSearch struct {
	ID        string `json:"id"`
	Nickname  string `json:"nickname"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Image     string `json:"image"`
}

func SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "missing query")
		return
	}

	query = strings.TrimSpace(query)

	rows, err := repository.Db.Query(`
		SELECT id, nickname, first_name, last_name, image
		FROM users
		WHERE nickname LIKE '%' || ? || '%'
		   OR first_name LIKE '%' || ? || '%'
		   OR last_name LIKE '%' || ? || '%'
		LIMIT 20
	`, query, query, query)

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "database error")
		return
	}
	defer rows.Close()

	var users []UserSearch
	for rows.Next() {
		var u UserSearch
		if err := rows.Scan(&u.ID, &u.Nickname, &u.FirstName, &u.LastName, &u.Image); err != nil {
			if err != sql.ErrNoRows {
				helper.RespondWithError(w, http.StatusInternalServerError, "scan error")
				return
			}
		}
		users = append(users, u)
	}

	json.NewEncoder(w).Encode(users)
}
