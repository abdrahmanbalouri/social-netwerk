package service

import (
	"database/sql"

	"social-network/app/repository/model"
	"social-network/app/utils"
)

func GetCommunFriends(userID string) ([]utils.User, int, error) {
	users, err := model.GetCommunFriends(userID)
	if err == sql.ErrNoRows {
		return nil, 200, nil
	} else if err != nil {
		return nil, 500, err
	}

	return users, 200, nil
}
