package service

import (
	"errors"
	"social-network/internal/repository/model"
)

func ToggleFollow(currentUserID, targetUserID string) (map[string]interface{}, error) {
	privacy, err := model.GetUserPrivacy(targetUserID)
	if err != nil {
		return nil, errors.New("failed to load user privacy")
	}

	isFollowed, err := model.IsFollowing(targetUserID, currentUserID)
	if err != nil {
		return nil, errors.New("database error")
	}

	if privacy == "private" {
		if isFollowed {
			model.Unfollow(targetUserID, currentUserID)
		} else {
			pending, _ := model.IsPending(targetUserID, currentUserID)
			if pending {
				model.CancelFollowRequest(targetUserID, currentUserID)
			} else {
				model.CreateFollowRequest(targetUserID, currentUserID)
			}
		}
	} else {
		if isFollowed {
			model.Unfollow(targetUserID, currentUserID)
		} else {
			model.Follow(targetUserID, currentUserID)
		}
	}

	followers := model.CountFollowers(targetUserID)
	following := model.CountFollowing(targetUserID)

	isPending, _ := model.IsPending(targetUserID, currentUserID)
	isFollowed, _ = model.IsFollowing(targetUserID, currentUserID)

	return map[string]interface{}{
		"followers":  followers,
		"following":  following,
		"isFollowed": isFollowed,
		"isPending":  isPending,
		"privacy":    privacy,
	}, nil
}
