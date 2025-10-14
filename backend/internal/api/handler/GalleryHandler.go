package handlers

import (
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func GalleryHandler(w http.ResponseWriter, r *http.Request) {
	q := `SELECT image_path, title FROM posts WHERE user_id = ? ORDER BY id DESC`
	userID := r.URL.Query().Get("id")
	rows, err := repository.Db.Query(q, userID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve gallery")
		return
	}
	defer rows.Close()

	var gallery []map[string]string
	for rows.Next() {
		var imagePath, title string
		if err := rows.Scan(&imagePath, &title); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to parse gallery data")
			return
		}

		gallery = append(gallery, map[string]string{
			"imagePath": imagePath,
			"title":     title,
		})
	}

	if err := rows.Err(); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve gallery")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, gallery)
}
