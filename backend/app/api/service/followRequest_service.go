package service

import (
	"database/sql"
	"errors"
	"net/http"

	"social-network/app/repository/model"
)

func GetFollowRequests(userID string) ([]map[string]interface{}, int, error) {
	// Get user's privacy settings
	privacy, err := model.GetUserPrivacy(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.New("user not found")
		}
		return nil, http.StatusInternalServerError, errors.New("failed to get user privacy")
	}

	// Clear follow requests if account is not private
	if privacy != "private" {
		err := model.ClearFollowRequests(userID)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.New("failed to clear follow requests")
		}
		// Return empty list for non-private accounts
		return []map[string]interface{}{}, http.StatusOK, nil
	}

	// Fetch follow requests for private account
	requests, err := model.FetchFollowRequests(userID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve follow requests")
	}

	return requests, http.StatusOK, nil
}