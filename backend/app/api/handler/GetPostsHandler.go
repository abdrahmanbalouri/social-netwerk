package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	service "social-network/app/api/service"
	"social-network/app/helper"
	"social-network/pkg/db/sqlite"
)

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
if r.Method !=  http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, " method not allowed ")
		return
}

	authUserID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	postID := parts[3]

	postResp, err := service.FetchPost(sqlite.Db, postID, authUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			helper.RespondWithError(w, http.StatusNotFound, "Post not found")
			return
		}
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch post")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, postResp)
}
