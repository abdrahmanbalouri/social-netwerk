package helper

import (
	"errors"

	"social-network/app/repository"
)

// CheckUserInGroup katchecki wach user kayn f group
func CheckUserInGroup(userID string, groupID string) error {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?)`
	err := repository.Db.QueryRow(query, userID, groupID).Scan(&exists)
	if err != nil {
		return err // error f query
	}
	if !exists {
		return errors.New("user is not a member of this group")
	}
	return nil
}
