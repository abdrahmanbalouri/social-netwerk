package handlers

import (
	"net/http"

	service "social-network/internal/api/service"
	"social-network/internal/helper"
)

func SearchUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query().Get("query")
	users, err := service.SearchUsers(query)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, users)
}
