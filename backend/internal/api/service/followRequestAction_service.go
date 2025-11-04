package service

import (
	"errors"
	"social-network/internal/repository/model"
)

func HandleFollowRequestAction(currentUserID, requesterID, action string) error {
	// Check request exists
	if !model.FollowRequestExists(currentUserID, requesterID) {
		return errors.New("Follow request not found")
	}

	// If user accepts → insert into followers
	if action == "accept" {
		if err := model.AddFollower(currentUserID, requesterID); err != nil {
			return err
		}
	}

	// In all cases → delete request
	return model.DeleteFollowRequest(currentUserID, requesterID)
}
