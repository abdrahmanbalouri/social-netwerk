package handlers

import (
	"net/http"
	"strings"

	"social-network/internal/api/service"
	"social-network/internal/helper"
)

func GetEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	UserId, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusNotFound, "Group not found")
		return
	}
	GrpID := parts[3]

	Events, err := service.GetEventsService(GrpID, UserId)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, Events)
}
