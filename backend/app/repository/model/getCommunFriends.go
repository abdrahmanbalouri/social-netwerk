package model

import (
	"social-network/app/utils"
	"social-network/pkg/db/sqlite"
)

func GetCommunFriends(userID string) ([]utils.User, error) {
	rows, err := sqlite.Db.Query(`
	SELECT DISTINCT u.id, u.first_name ,u.last_name, u.image
	FROM users u
	INNER JOIN followers f 
	ON (u.id = f.follower_id OR u.id = f.user_id)
	WHERE (f.user_id = ? OR f.follower_id = ?);
	`, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	friends := []utils.User{}

	for rows.Next() {
		friend := utils.User{}

		if err := rows.Scan(&friend.ID, &friend.First_name, &friend.Last_name, &friend.Image); err != nil {
			return nil, err
		}
		if userID != friend.ID {
			friends = append(friends, friend)
		}
	}

	return friends, nil
}
