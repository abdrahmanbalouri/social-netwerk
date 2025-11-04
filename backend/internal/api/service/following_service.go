package service

import (
	"errors"
	"social-network/internal/repository/model"
)

func GetFollowing(currentUserID, targetUserID string) ([]map[string]interface{}, error) {
	// 1) Verify user exists
	if !model.UserExists(targetUserID) {
		return nil, errors.New("User not found")
	}

	// 2) Check privacy rules
	privacy, isFollowing, err := model.GetPrivacyAndFollowing(currentUserID, targetUserID)
	if err != nil {

	}
	if privacy == "private" && isFollowing == 0 && currentUserID != targetUserID {
		return nil, errors.New("This account is private")
	}

	// 3) Get following list
	return model.GetFollowingList(targetUserID)
}
