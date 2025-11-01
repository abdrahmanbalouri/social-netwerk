package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

type LikeRequest struct {
	PostID string `json:"postId"`
}

func LikesGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Authenticate user
	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var newLike LikeRequest
	if err := json.NewDecoder(r.Body).Decode(&newLike); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request format1")
		return
	}

	// check if the user is a member of the group
	query := `SELECT p.group_id, EXISTS(SELECT 1 FROM group_members gm WHERE gm.user_id = ? AND gm.group_id = p.group_id) FROM posts p WHERE p.id = ?;`
	var grpID string
	var isMember bool
	err = repository.Db.QueryRow(query, userID, newLike.PostID).Scan(&grpID, &isMember)
	if err != nil {
		if err == sql.ErrNoRows {
			helper.RespondWithError(w, http.StatusNotFound, "Failed to get the group's id or post not exist")
		}
		helper.RespondWithError(w, http.StatusInternalServerError, "error finding the group/post")
		return
	}
	if !isMember {
		helper.RespondWithError(w, http.StatusUnauthorized, "User is not a member of the group")
		return
	}

	// Check if the user has already liked the post
	var existingLikeID string
	err = repository.Db.QueryRow(`
		SELECT id FROM likes 
		WHERE user_id = ? AND liked_item_id = ? AND liked_item_type = 'post'
	`, userID, newLike.PostID).Scan(&existingLikeID)

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
	likeID := helper.GenerateUUID()
	_, err = repository.Db.Exec(`
		INSERT INTO likes (id, user_id, liked_item_id, liked_item_type, created_at)
		VALUES (?, ?, ?, 'post', ?)
	`, likeID, userID, newLike.PostID, time.Now())
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to add like")
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Like added"})
}
