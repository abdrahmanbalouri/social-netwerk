package handlers

import (
	"net/http"

	service "social-network/internal/api/sevice"
	"social-network/internal/helper"
)

func GalleryHandler(w http.ResponseWriter, r *http.Request) {
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
