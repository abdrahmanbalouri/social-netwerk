package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func LikesGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get postID and groupID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		helper.RespondWithError(w, http.StatusBadRequest, "Post ID and Group ID are required")
		return
	}
	postID := pathParts[3]
	groupID := pathParts[4]

	// check if the user is a member of the group
	query := `
	SELECT EXISTS(
		SELECT 1 FROM group_members gm 
		WHERE gm.user_id = ? AND gm.group_id = ?
	)
	`
	var isMember bool
	err = repository.Db.QueryRow(query, userID, groupID).Scan(&isMember)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	if !isMember {
		helper.RespondWithError(w, http.StatusUnauthorized, "User is not a member of the group")
		return
	}

	// check if the post exists in the group
	var exists bool
	err = repository.Db.QueryRow(`SELECT EXISTS(SELECT 1 FROM group_posts WHERE id = ? AND group_id = ?)`, postID, groupID).Scan(&exists)
	if err != nil || !exists {
		helper.RespondWithError(w, http.StatusNotFound, "Post not found in this group")
		return
	}

	// Check if user already liked the post
	var existingLikeID string
	err = repository.Db.QueryRow(`
		SELECT id FROM likesgroups
		WHERE user_id = ? AND liked_item_id = ? AND liked_item_type = 'post'
	`, userID, postID).Scan(&existingLikeID)

	if err == nil {
		// Unlike
		_, err = repository.Db.Exec(`DELETE FROM likesgroups WHERE id = ?`, existingLikeID)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to remove like")
			return
		}
		helper.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Like removed"})
		return
	} else if err != sql.ErrNoRows {
		helper.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Add like
	likeID := helper.GenerateUUID()
	_, err = repository.Db.Exec(`
		INSERT INTO likesgroups (id, user_id, liked_item_id, liked_item_type, created_at)
		VALUES (?, ?, ?, 'post', ?)
	`, likeID, userID, postID, time.Now())
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to add like")
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Like added"})
}
