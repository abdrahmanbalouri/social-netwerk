package handlers

import (
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository/post"
)

func AllpostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	posts, err := post.GetAllPosts()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve posts")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, posts)
}
