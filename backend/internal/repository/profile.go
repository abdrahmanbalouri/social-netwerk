package repository

import "database/sql"

type Profile struct {
	UserID      string
	DisplayName string
	Privacy     string
	Cover       string
	Avatar      string
}

// GetProfileFiles retrieves existing cover and avatar filenames
func GetProfileFiles(db *sql.DB, userID string) (cover, avatar string, err error) {
	err = db.QueryRow("SELECT cover, image FROM users WHERE id = ?", userID).Scan(&cover, &avatar)
	return
}

// UpdateProfile updates profile information in the database
func UpdateProfile(db *sql.DB, profile Profile) error {
	_, err := db.Exec(`
		UPDATE users
		SET about = ?, privacy = ?, cover = ?, image = ?
		WHERE id = ?`,
		profile.DisplayName,
		profile.Privacy,
		profile.Cover,
		profile.Avatar,
		profile.UserID,
	)
	return err
}
