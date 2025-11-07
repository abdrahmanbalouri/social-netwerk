package handlers

import (
	"net/http"

	service "social-network/internal/api/service"
	"social-network/internal/helper"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := service.ValidateSession(c.Value)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	resp := struct {
		Message string `json:"message"`
		UserID  string `json:"user_id"`
	}{
		Message: "authorized",
		UserID:  userID,
	}

	helper.RespondWithJSON(w, http.StatusOK, resp)
}
