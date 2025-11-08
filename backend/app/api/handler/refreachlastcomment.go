package handlers

import (
	"fmt"
	"net/http"
	"strings"

	service "social-network/app/api/service"
	"social-network/app/helper"
)

func GetlastcommnetGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 { // /api/comments/<commentID>/<groupID>
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	commentID := parts[3]
	groupID := parts[4]

	commentMap, err := service.GetLastCommentGroup(commentID, userID, groupID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get comment: %v", err))
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, commentMap)
}
