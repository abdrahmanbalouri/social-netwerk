package model

import "social-network/pkg/db/sqlite"

func GetFollowingList(targetUserID string) ([]map[string]interface{}, error) {
	rows, err := sqlite.Db.Query(`
		SELECT u.id, u.first_name, u.last_name, u.image 
		FROM followers f
		JOIN users u ON u.id = f.user_id
		WHERE f.follower_id = ?;
	`, targetUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []map[string]interface{}
	for rows.Next() {
		var idU, first, last, img string
		if err := rows.Scan(&idU, &first, &last, &img); err != nil {
			return nil, err
		}
		following = append(following, map[string]interface{}{
			"id":         idU,
			"first_name": first,
			"last_name":  last,
			"image":      img,
		})
	}
	return following, nil
}
