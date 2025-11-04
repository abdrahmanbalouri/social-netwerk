package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	service "social-network/internal/api/service"
	"social-network/internal/helper"
)

func Editor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	displayName := r.FormValue("displayName")
	privacy := r.FormValue("privacy")

	var coverFile, avatarFile io.ReadCloser
	var coverName, avatarName string

	cFile, cHeader, cErr := r.FormFile("cover")
	if cErr == nil {
		coverFile = cFile
		coverName = cHeader.Filename
	}

	aFile, aHeader, aErr := r.FormFile("avatar")
	if aErr == nil {
		avatarFile = aFile
		avatarName = aHeader.Filename
	}

	coverFilename, avatarFilename, err := service.UpdateUserProfile(userID, displayName, privacy, coverFile, coverName, avatarFile, avatarName)
	if err != nil {
		http.Error(w, "Failed to update profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"message": "Profile updated successfully",
		"cover":   coverFilename,
		"avatar":  avatarFilename,
	})
}
