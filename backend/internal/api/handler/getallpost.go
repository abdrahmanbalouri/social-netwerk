package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"social-network/internal/helper"
	"social-network/internal/repository/post"
)

func AllpostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
	id := r.URL.Query().Get("userId")
	if id == "0" {
		userID, _ := helper.AuthenticateUser(r)
		id = strconv.Itoa(userID)
	}
	fmt.Println("User ID:", id)

	posts, err := post.GetAllPosts(id)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve posts")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, posts)
}
