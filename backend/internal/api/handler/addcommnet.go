package handlers

import (
	"fmt"
	"net/http"

	service "social-network/internal/api/sevice"
	"social-network/internal/helper"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = r.ParseMultipartForm(20 << 20)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	postID := r.FormValue("postId")
	content := r.FormValue("content")
	whatis := r.FormValue("whatis")
	groupID := r.FormValue("groupId")

	// Extract media file if exists
	var mediaFileHeader map[string]interface{}
	file, header, err := r.FormFile("media")
	if err == nil {
		defer file.Close()
		mediaFileHeader = map[string]interface{}{
			"file":     file,
			"filename": header.Filename,
			"size":     header.Size,
		}
	}
	fmt.Println(mediaFileHeader["size"])

	commentID, mediaPath, err := service.CreateComment(userID, postID, content, whatis, groupID, mediaFileHeader)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"message":    "Comment created successfully",
		"comment_id": commentID,
		"media":      mediaPath,
	})
}
