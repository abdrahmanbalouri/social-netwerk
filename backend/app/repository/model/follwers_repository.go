package model

import "social-network/pkg/db/sqlite"

func UserExists(id string) bool {
	var tmp string
	err := sqlite.Db.QueryRow(`SELECT id FROM users WHERE id = ?`, id).Scan(&tmp)
	return err == nil
}

func GetPrivacyAndFollowing(currentUserID, targetUserID string) (privacy string, isFollowing int, err error) {
	err = sqlite.Db.QueryRow(`
SELECT 
    u.privacy,
    CASE WHEN f.follower_id IS NOT NULL THEN 1 ELSE 0 END
FROM users u
LEFT JOIN followers f 
    ON f.user_id = u.id 
   AND f.follower_id = ?
WHERE u.id = ?;
`, currentUserID, targetUserID).Scan(&privacy, &isFollowing)

	return
}

func GetFollowersList(targetUserID string) ([]map[string]interface{}, error) {
	rows, err := sqlite.Db.Query(`
	SELECT u.id , u.first_name, u.last_name, u.image
	FROM followers f
	JOIN users u ON u.id = f.follower_id
	WHERE f.user_id = ?;`, targetUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []map[string]interface{}
	for rows.Next() {
		var idU, first, last, img string
		if err := rows.Scan(&idU, &first, &last, &img); err != nil {
			return nil, err
		}
		followers = append(followers, map[string]interface{}{
			"id":         idU,
			"first_name": first,
			"last_name":  last,
			"image":      img,
		})
	}
	return followers, nil
}
