package service

import (
	"errors"

	"social-network/app/repository/model"
)

func GetFollowRequests(userID string) ([]map[string]interface{}, error) {
	privacy, err := model.GetUserPrivacy(userID)
	if err != nil {
		return nil, errors.New("failed to get user privacy")
	}

	if privacy != "private" {
		err := model.ClearFollowRequests(userID)
		if err != nil {
			return nil, errors.New("failed to clear follow requests for non-private account")
		}
	}
	return model.FetchFollowRequests(userID)
}
