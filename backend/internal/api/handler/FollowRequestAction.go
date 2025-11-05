package handlers

import (
	"encoding/json"
	"net/http"

	"social-network/internal/api/service"
	"social-network/internal/helper"
)

func FollowRequestAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, " method not allowed ")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	var req struct {
		ID     string `json:"id"`
		Action string `json:"action"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	err = service.HandleFollowRequestAction(userID, req.ID, req.Action)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]string{
		"status": "success",
		"action": req.Action,
	})
}
