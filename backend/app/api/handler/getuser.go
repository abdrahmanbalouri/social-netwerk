package handlers

import (
	"net/http"

	service "social-network/app/api/service"
	"social-network/app/helper"
	"social-network/app/repository"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	currentUserID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	users, err := service.GetUsers(repository.Db, currentUserID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, users)
}
