package model

import (
	"database/sql"
)

type User struct {
	ID         string `json:"id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Image      string `json:"image"`
}

func SearchUsersInDB(db *sql.DB, search string) ([]map[string]any, error) {
	rows, err := db.Query(`
        SELECT id, nickname, first_name, last_name, image 
        FROM users 
        WHERE lower(nickname) LIKE ? OR lower(first_name) LIKE ? 
		OR lower(last_name) LIKE ? 
		OR lower(CONCAT(first_name,' ',last_name)) LIKE ?
        LIMIT 10
    `, search, search, search, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]any
	for rows.Next() {
		var id, nickname, firstName, lastName, image string
		if err := rows.Scan(&id, &nickname, &firstName, &lastName, &image); err != nil {
			continue
		}
		users = append(users, map[string]any{
			"id":         id,
			"nickname":   nickname,
			"first_name": firstName,
			"last_name":  lastName,
			"image":      image,
		})
	}
	return users, nil
}

func FetchAllUsers(db *sql.DB, excludeUserID string) ([]User, error) {
	rows, err := db.Query("SELECT id, first_name, last_name, image FROM users WHERE id != ?", excludeUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.First_name, &u.Last_name, &u.Image); err != nil {
			continue
		}
		users = append(users, u)
	}
	return users, nil
}
