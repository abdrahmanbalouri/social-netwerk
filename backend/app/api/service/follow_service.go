package service

import (
	"database/sql"
	"errors"
	"net/http"

	"social-network/app/repository/model"
)

func ToggleFollow(currentUserID, targetUserID string) (map[string]interface{}, int, error) {
	// Validate: can't follow yourself
	if currentUserID == targetUserID {
		return nil, http.StatusBadRequest, errors.New("you can't follow yourself")
	}

	// Get target user's privacy settings
	privacy, err := model.GetUserPrivacy(targetUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.New("user not found")
		}
		return nil, http.StatusInternalServerError, errors.New("failed to load user privacy")
	}

	// Check if already following
	isFollowed, err := model.IsFollowing(targetUserID, currentUserID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("database error while checking follow status")
	}

	// Handle private accounts
	if privacy == "private" {
		if isFollowed {
			// Unfollow
			err = model.Unfollow(targetUserID, currentUserID)
			if err != nil {
				return nil, http.StatusInternalServerError, errors.New("failed to unfollow user")
			}
		} else {
			// Check if request is pending
			pending, err := model.IsPending(targetUserID, currentUserID)
			if err != nil {
				return nil, http.StatusInternalServerError, errors.New("database error while checking pending request")
			}

			if pending {
				// Cancel pending request
				err = model.CancelFollowRequest(targetUserID, currentUserID)
				if err != nil {
					return nil, http.StatusInternalServerError, errors.New("failed to cancel follow request")
				}
			} else {
				// Create new follow request
				err = model.CreateFollowRequest(targetUserID, currentUserID)
				if err != nil {
					return nil, http.StatusInternalServerError, errors.New("failed to create follow request")
				}
			}
		}
	} else {
		// Handle public accounts
		if isFollowed {
			// Unfollow
			err = model.Unfollow(targetUserID, currentUserID)
			if err != nil {
				return nil, http.StatusInternalServerError, errors.New("failed to unfollow user")
			}
		} else {
			// Follow directly
			err = model.Follow(targetUserID, currentUserID)
			if err != nil {
				return nil, http.StatusInternalServerError, errors.New("failed to follow user")
			}
		}
	}

	// Get updated stats
	followers := model.CountFollowers(targetUserID)
	following := model.CountFollowing(targetUserID)

	// Get updated follow status
	isPending, err := model.IsPending(targetUserID, currentUserID)
	if err != nil {
		// Non-critical error, just log it
		isPending = false
	}

	isFollowed, err = model.IsFollowing(targetUserID, currentUserID)
	if err != nil {
		// Non-critical error, just log it
		isFollowed = false
	}

	return map[string]interface{}{
		"followers":  followers,
		"following":  following,
		"isFollowed": isFollowed,
		"isPending":  isPending,
		"privacy":    privacy,
	}, http.StatusOK, nil
}