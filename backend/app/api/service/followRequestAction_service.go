package service

import (
	"database/sql"
	"errors"
	"net/http"

	"social-network/app/repository/model"
)

func HandleFollowRequestAction(currentUserID, requesterID, action string) (int, error) {
	// Validate action
	if action != "accept" && action != "reject" {
		return http.StatusBadRequest, errors.New("invalid action: must be 'accept' or 'reject'")
	}

	// Validate: can't be the same user
	if currentUserID == requesterID {
		return http.StatusBadRequest, errors.New("invalid request")
	}

	// Get current user's privacy settings
	privacy, err := model.GetUserPrivacy(currentUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, errors.New("user not found")
		}
		return http.StatusInternalServerError, errors.New("failed to get user privacy")
	}

	// Only private accounts can manage follow requests
	if privacy != "private" {
		err := model.ClearFollowRequests(currentUserID)
		if err != nil {
			return http.StatusInternalServerError, errors.New("failed to clear follow requests")
		}
		return http.StatusBadRequest, errors.New("follow requests can only be managed for private accounts")
	}

	// Check if follow request exists
	exists := model.FollowRequestExists(currentUserID, requesterID)
	if !exists {
		return http.StatusNotFound, errors.New("follow request not found")
	}

	// If user accepts → insert into followers
	if action == "accept" {
		err := model.AddFollower(currentUserID, requesterID)
		if err != nil {
			return http.StatusInternalServerError, errors.New("failed to accept follow request")
		}
	}

	// In all cases (accept or reject) → delete the request
	err = model.DeleteFollowRequest(currentUserID, requesterID)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to remove follow request")
	}

	return http.StatusOK, nil
}