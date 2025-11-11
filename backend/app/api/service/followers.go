package service

import (
	"social-network/app/repository/model"
	"social-network/app/utils"
	"social-network/pkg/db/sqlite"
)

// GetFollowersService handles the business logic
func GetFollowersService(useID string) ([]utils.User, error) {
	// 1. Authenticate user

	// 2. Fetch followers from DB
	followers, err := model.GetFollowersByUser(sqlite.Db, useID)
	if err != nil {
		return nil, err
	}

	// 3. Optional: extra logic (filter, sort, etc.)
	return followers, nil
}
