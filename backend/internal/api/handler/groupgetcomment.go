package handlers

import (
	"net/http"
	"strconv"
	"strings"

	service "social-network/internal/api/service"
	"social-network/internal/helper"
)

func GetCommentsGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	postID := parts[3]
	offsetStr := parts[4]
	groupID := parts[5]

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	comments, err := service.GetCommentsGroup(userID, groupID, postID, offset)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, comments)
}
