package handlers

import (
	"io"
	"net/http"
	"strings"

	service "social-network/internal/api/sevice"
	"social-network/internal/helper"
	"social-network/internal/repository/middleware"
)

func Createpost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	userID, err := middleware.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	title := helper.Skip(strings.TrimSpace(r.FormValue("title")))
	content := helper.Skip(strings.TrimSpace(r.FormValue("content")))
	visibility := strings.TrimSpace(r.FormValue("visibility"))
	allowedUsers := strings.TrimSpace(r.FormValue("allowed_users"))

	var fileHeader io.ReadCloser
	var filename string
	file, header, err := r.FormFile("image")
	if err == nil {
		fileHeader = file
		filename = header.Filename
	}

	postID, err := service.CreatePost(userID, title, content, visibility, allowedUsers, fileHeader, filename)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "Post created successfully",
		"post_id": postID,
	})
}
