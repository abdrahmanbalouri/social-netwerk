package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"social-network/internal/helper"
	"social-network/internal/repository/post"
)

func AllpostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 { // [0]= "", [1]=api, [2]=Getcomments, [3]=postID, [4]=offset
		helper.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	offsetStr := parts[3]
	userId := r.URL.Query().Get("userId")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	d, err := helper.AuthenticateUser(r)
	if err != nil {
		fmt.Println(d)
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	posts, err := post.GetAllPosts(userId, r, offset)
	if posts == nil {
		helper.RespondWithJSON(w, http.StatusOK, posts)
		return
	}
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve posts")
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, posts)
}
