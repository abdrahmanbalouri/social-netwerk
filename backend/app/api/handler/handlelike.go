package handlers

import (
	"net/http"
	"strings"

	service "social-network/app/api/service"
	"social-network/app/helper"
	"social-network/app/repository"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract post ID
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusBadRequest, "Post ID is required")
		return
	}
	postID := parts[3]

	// Authenticate user
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	action, err := service.TogglePostLike(repository.Db, userID, postID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to toggle like")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Like " + action})
}
