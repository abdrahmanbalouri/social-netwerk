package service

import (
	"errors"
	"fmt"

	"social-network/internal/repository/model"
)

func GetFollowers(currentUserID, targetUserID string) ([]map[string]interface{}, error) {
	// 1) Verify user exists
	if !model.UserExists(targetUserID) {
		fmt.Println("dededede")
		return nil, errors.New("User not found")
	}

	// 2) Check privacy rules
	privacy, isFollowing := model.GetPrivacyAndFollowing(currentUserID, targetUserID)
	if privacy == "private" && isFollowing == 0 && currentUserID != targetUserID {
		return nil, errors.New("This account is private")
	}

	// 3) Get followers list
	return model.GetFollowersList(targetUserID)
}
