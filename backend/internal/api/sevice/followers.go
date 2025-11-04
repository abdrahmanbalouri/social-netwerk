package service

import (
	"social-network/internal/repository"
	"social-network/internal/utils"
)

// GetFollowersService handles the business logic
func GetFollowersService(useID string) ([]utils.User, error) {
	// 1. Authenticate user

	// 2. Fetch followers from DB
	followers, err := repository.GetFollowersByUser(repository.Db, useID)
	if err != nil {
		return nil, err
	}

	// 3. Optional: extra logic (filter, sort, etc.)
	return followers, nil
}
