package model

import "social-network/app/repository"

func FetchFollowRequests(userID string) ([]map[string]interface{}, error) {
	rows, err := repository.Db.Query(`
		SELECT u.id, u.first_name, u.last_name, u.image
		FROM users u
		JOIN follow_requests f ON f.follower_id = u.id
		WHERE f.user_id = ?;
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []map[string]interface{}
	for rows.Next() {
		var id, first, last, img string
		if err := rows.Scan(&id, &first, &last, &img); err != nil {
			return nil, err
		}

		requests = append(requests, map[string]interface{}{
			"id":         id,
			"first_name": first,
			"last_name":  last,
			"image":      img,
		})
	}

	return requests, nil
}

func ClearFollowRequests(userID string) error {
	_, err := repository.Db.Exec(`
		DELETE FROM follow_requests
		WHERE user_id = ?;
	`, userID)
	return err
}
