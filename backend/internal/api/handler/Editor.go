package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

type ProfilePayload struct {
	DisplayName string `json:"displayName"`
	Privacy     string `json:"privacy"`
	Cover       string `json:"cover"`
	Avatar      string `json:"avatar"`
}

func Editor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	userid, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	displayName := r.FormValue("displayName")
	privacy := r.FormValue("privacy")

	coverFile, coverHeader, err := r.FormFile("cover")
	var coverFilename string
	if err == nil {
		defer coverFile.Close()
		coverPath := "../frontend/my-app/public/uploads/" + coverHeader.Filename
		out, _ := os.Create(coverPath)
		defer out.Close()
		io.Copy(out, coverFile)
		coverFilename = coverHeader.Filename
	} else {
		coverFilename = "" // أو خليه بنفس القيمة القديمة فالـ DB
	}

	avatarFile, avatarHeader, err := r.FormFile("avatar")
	var avatarFilename string
	if err == nil {
		defer avatarFile.Close()
		avatarPath := "../frontend/my-app/public/" + avatarHeader.Filename
		out, _ := os.Create(avatarPath)
		defer out.Close()
		io.Copy(out, avatarFile)
		avatarFilename = avatarHeader.Filename
	} else {
		avatarFilename = "default.png"
	}

	// Update DB
	_, err = repository.Db.Exec(`
	UPDATE users
	SET about = ?, privacy = ?, cover = ?, image = ?
	WHERE id = ?`,
		displayName,
		privacy,
		coverFilename,
		avatarFilename,
		userid,
	)

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
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
