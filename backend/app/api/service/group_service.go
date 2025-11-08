package service

import (
	"fmt"
	"time"

	"social-network/app/helper"
	"social-network/app/repository"
	"social-network/app/repository/model"
	"social-network/app/utils"
)

func GetUserGroups(userID string) ([]utils.Group, error) {
	return model.FetchUserGroups(userID)
}

func GetAllAvailableGroups(userID string) ([]utils.Group, error) {
	return model.FetchAllAvailableGroups(userID)
}

func CreateNewGroup(adminID string, newGroup utils.GroupRequest) (utils.Group, error) {
	tx, err := repository.Db.Begin()
	if err != nil {
		return utils.Group{}, fmt.Errorf("Failed to start database transaction")
	}
	defer tx.Rollback()

	groupID := helper.GenerateUUID()

	if err := model.InsertNewGroup(tx, groupID.String(), newGroup.Title, newGroup.Description, adminID); err != nil {
		return utils.Group{}, fmt.Errorf("Failed to create new group")
	}

	if err := model.InsertAdminAsMember(tx, adminID, groupID.String()); err != nil {
		return utils.Group{}, fmt.Errorf("Failed to insert admin into group members table")
	}

	for _, userID := range newGroup.InvitedUsers {
		rowID := helper.GenerateUUID()
		createdAt := time.Now().UTC()
		if err := model.InsertGroupInvitation(tx, rowID.String(), groupID.String(), userID, adminID, createdAt); err != nil {
			fmt.Println("Failed to insert invited user into group_invitation table :", err)
			return utils.Group{}, fmt.Errorf("Failed to insert invited user into group_invitation table")
		}
	}

	if err := tx.Commit(); err != nil {
		return utils.Group{}, fmt.Errorf("Failed to commit transaction")
	}

	createdGroup, err := model.FetchCreatedGroup(groupID.String())
	if err != nil {
		fmt.Println("Failed to fetch created group:", err)
		return utils.Group{}, fmt.Errorf("Group created but failed to fetch it")
	}

	return createdGroup, nil
}

func HandleGroupInvitation(groupID, userID string, newInvitation utils.GroupInvitation) (map[string]any, error) {
	tx, err := repository.Db.Begin()
	if err != nil {
		return nil, fmt.Errorf("Failed to start database transaction")
	}
	defer tx.Rollback()

	invitationID := helper.GenerateUUID()

	if newInvitation.InvitationType == "join" {
		// Check if similar invitation exists
		exists, err := model.CheckExistingInvitation(tx, userID, groupID)
		if err != nil {
			return nil, fmt.Errorf("Database error: %v", err)
		}
		if exists {
			return nil, fmt.Errorf("There is another invitation with the same credentials")
		}

		// Check membership
		isMember, err := model.CheckGroupMembership(tx, userID, groupID)
		if err != nil {
			return nil, fmt.Errorf("Failed to check group membership")
		}
		if isMember {
			return nil, fmt.Errorf("You are already a member of this group")
		}

		if err := model.InsertJoinRequest(tx, invitationID.String(), groupID, userID); err != nil {
			return nil, fmt.Errorf("Error sending the invitation: %v", err)
		}

	} else {
		for _, invitedUser := range newInvitation.InvitedUsers {
			exists, err := model.CheckUserMembershipOrInvitation(tx, invitedUser, groupID)
			if err != nil {
				return nil, fmt.Errorf("Error checking for existing membership or invitation")
			}
			if exists {
				continue
			}

			userExists, err := model.CheckUserExists(tx, invitedUser)
			if err != nil {
				return nil, fmt.Errorf("Database error")
			}
			if !userExists {
				return nil, fmt.Errorf("The invited user isn't registered")
			}

			if err := model.InsertInvitation(tx, invitationID.String(), groupID, invitedUser, userID); err != nil {
				return nil, fmt.Errorf("Error sending the invitation: %v", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("Failed to commit transaction")
	}

	return map[string]any{
		"invitation_id": invitationID,
		"message":       "Invitation successfully processed",
	}, nil
}

func ProcessGroupInvitationResponse(userID string, response utils.GroupResponse) error {
	fmt.Println("USER ID IS :", userID)
	fmt.Println("RESPONSE IS :", response)
	// Determine the actual user ID if the response is not from an invited user
	if response.InvitationType != "invitation" {
		dbUserID, err := model.GetUserIDByInvitation(response.InvitationID)
		if err != nil {
			return fmt.Errorf("failed to fetch user ID: %v", err)
		}
		userID = dbUserID
	}

	// Fetch group ID related to this invitation
	groupID, err := model.GetGroupIDByInvitation(response.InvitationID, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch group ID: %v", err)
	}

	tx, err := repository.Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start database transaction")
	}
	defer tx.Rollback()

	// If user accepted, add to members
	if response.Response == "accept" {
		if err := model.AddUserToGroup(tx, userID, groupID); err != nil {
			return fmt.Errorf("failed to add user to group: %v", err)
		}
	}

	// Delete invitation after response
	if err := model.DeleteInvitation(tx, response.InvitationID); err != nil {
		return fmt.Errorf("failed to delete invitation: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
