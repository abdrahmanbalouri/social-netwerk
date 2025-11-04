package model

import (
	"database/sql"
	"social-network/internal/utils"
)

// GetFollowersByUser fetches all followers of a given user
func GetFollowersByUser(db *sql.DB, userID string) ([]utils.User, error) {
	query := `
		SELECT u.id, u.first_name, u.last_name, u.image
		FROM followers f
		JOIN users u ON u.id = f.follower_id
		WHERE f.user_id = ?;
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []utils.User
	for rows.Next() {
		var u utils.User
		if err := rows.Scan(&u.ID, &u.First_name, &u.Last_name, &u.Image); err != nil {
			return nil, err
		}
		followers = append(followers, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}
