package repository

import "database/sql"

// PostGallery holds minimal post info for gallery
type PostGallery struct {
	ImagePath string
	Title     string
}

// GetUserGallery fetches posts for a user
func GetUserGallery(db *sql.DB, userID string) ([]PostGallery, error) {
query := `SELECT image_path, title FROM posts WHERE user_id = ? AND image_path != "" ORDER BY id DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gallery []PostGallery
	for rows.Next() {
		var pg PostGallery
		if err := rows.Scan(&pg.ImagePath, &pg.Title); err != nil {
			return nil, err
		}
		gallery = append(gallery, pg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return gallery, nil
}
