package service

import (
	"database/sql"
	"errors"
	"net/http"

	"social-network/app/repository/model"
)

func GetFollowers(currentUserID, targetUserID string) ([]map[string]interface{}, int, error) {
	// 1) Verify user exists
	exists := model.UserExists(targetUserID)
	if !exists {
		return nil, http.StatusNotFound, errors.New("user not found")
	}

	// 2) Check privacy rules
	privacy, isFollowing, err := model.GetPrivacyAndFollowing(currentUserID, targetUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.New("user not found")
		}
		return nil, http.StatusInternalServerError, errors.New("failed to check privacy settings")
	}

	// 3) Check if user has permission to view followers
	if privacy == "private" && isFollowing == 0 && currentUserID != targetUserID {
		return nil, http.StatusForbidden, errors.New("this account is private")
	}

	// 4) Get followers list
	followers, err := model.GetFollowersList(targetUserID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve followers list")
	}

	return followers, http.StatusOK, nil
}