package handlers

import (
	"io"
	"net/http"
	"strings"
	"time"

	service "social-network/internal/api/service"
	"social-network/internal/helper"
)

func CreateStories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	content := r.FormValue("content")
	bgColor := r.FormValue("bg_color")
	content = helper.Skip(strings.TrimSpace(content))

	var file io.ReadCloser
	var filename string
	imgFile, imgHeader, imgErr := r.FormFile("image")
	if imgErr == nil {
		file = imgFile
		filename = imgHeader.Filename
	}

	imagePath, err := service.CreateStory(userID, content, bgColor, file, filename)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"success":    true,
		"image_url":  imagePath,
		"created_at": time.Now().Format("2006-01-02 15:04:05"),
	})
}
