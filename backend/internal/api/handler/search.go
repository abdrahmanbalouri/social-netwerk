package handlers

import (
	"net/http"
	"strings"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func SearchUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		helper.RespondWithJSON(w, http.StatusOK, []any{})
		return
	}

	search := "%" + strings.ToLower(query) + "%"

	rows, err := repository.Db.Query(`
		SELECT id, nickname, first_name, last_name, image 
		FROM users 
		WHERE lower(nickname) LIKE ? OR lower(first_name) LIKE ? OR lower(last_name) LIKE ?
		LIMIT 10
	`, search, search, search)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var users []map[string]any
	for rows.Next() {
		var id, nickname, firstName, lastName, image string
		err = rows.Scan(&id, &nickname, &firstName, &lastName, &image)
		if err != nil {
			continue
		}
		users = append(users, map[string]any{
			"id":         id,
			"nickname":   nickname,
			"first_name": firstName,
			"last_name":  lastName,
			"image":      image,
		})
	}

	helper.RespondWithJSON(w, http.StatusOK, users)
}
