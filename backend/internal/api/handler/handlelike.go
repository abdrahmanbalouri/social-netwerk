package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"social-network/internal/helper"
	"social-network/internal/repository"

	"github.com/google/uuid"
)

// LikeHandler handles liking or unliking a post
func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract post ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		helper.RespondWithError(w, http.StatusBadRequest, "Post ID is required")
		return
	}
	postID := pathParts[3]

	// Authenticate user
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
   
	


	var existingLikeID string
	err = repository.Db.QueryRow(`
		SELECT id FROM likes 
		WHERE user_id = ? AND liked_item_id = ? AND liked_item_type = 'post'
	`, userID, postID).Scan(&existingLikeID)

	if err == nil {
		// Like exists, so remove it (unlike)
		_, err = repository.Db.Exec(`
			DELETE FROM likes 
			WHERE id = ?
		`, existingLikeID)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to remove like")
			return
		}
		helper.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Like removed"})
		return
	} else if err != sql.ErrNoRows {
		// Unexpected database error
		helper.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// No existing like, so add a new like
	likeID := uuid.New().String()
	_, err = repository.Db.Exec(`
		INSERT INTO likes (id, user_id, liked_item_id, liked_item_type, created_at)
		VALUES (?, ?, ?, 'post', ?)
	`, likeID, userID, postID, time.Now())
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to add like")
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Like added"})
}
