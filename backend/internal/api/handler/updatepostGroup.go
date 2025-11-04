package handlers

import (
	"net/http"
	"strings"

	service "social-network/internal/api/sevice"
	"social-network/internal/helper"
)

func GetGroupPostByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	currentUserID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 { // /groups/{groupID}/posts/{postID}
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	postID := parts[3]
	groupID := parts[4]

	postMap, err := service.GetGroupPost(postID, currentUserID, groupID)
	if err != nil {
		helper.RespondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, postMap)
}
