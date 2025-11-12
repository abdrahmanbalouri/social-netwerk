package handlers

import (
	"net/http"

	service "social-network/app/api/service"
	"social-network/app/helper"
)

func GalleryHandler(w http.ResponseWriter, r *http.Request) {

if r.Method != http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "methode not allowed")
		return
}

	userID := r.URL.Query().Get("id")
	if userID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	gallery, err := service.FetchUserGallery(userID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve gallery")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, gallery)
}
