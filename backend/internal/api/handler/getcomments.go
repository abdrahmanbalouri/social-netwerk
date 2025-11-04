package handlers

import (
	"net/http"
	"strconv"
	"strings"

	service "social-network/internal/api/sevice"
	"social-network/internal/helper"
)

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		helper.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	postID := parts[3]
	offset, err := strconv.Atoi(parts[4])
	if err != nil || offset < 0 {
		offset = 0
	}

	UserId, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
      ok,_  := helper.Canshowdata(UserId,postID)
	if !ok {
    helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch comments")
		return
	}
	

	comments, err := service.FetchComments(postID, offset)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch comments")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, comments)
}
