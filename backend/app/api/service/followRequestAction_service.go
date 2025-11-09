package service

import (
	"errors"

	"social-network/app/repository/model"
)

func HandleFollowRequestAction(currentUserID, requesterID, action string) error {
	// Check request exists

	privacy, err := model.GetUserPrivacy(currentUserID)
	if err != nil {
		return errors.New("failed to get user privacy")
	}

	if privacy != "private" {
		err := model.ClearFollowRequests(currentUserID)
		if err != nil {
			return errors.New("failed to clear follow requests for non-private account")
		}
		return errors.New("follow requests can only be managed for private accounts")
	}

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
