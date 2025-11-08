package model

import (
	"database/sql"
	"time"
)

type Story struct {
	UserID   string
	Content  string
	ImageURL string
	BgColor  string
}
type Storyapi struct {
	ID        string
	UserID    string
	Content   string
	ImageURL  string
	BGColor   string
	CreatedAt time.Time
	ExpiresAt sql.NullTime
	FirstName string
	LastName  string
	Profile   string
}

// InsertStory inserts a story into the database
func InsertStory(db *sql.DB, story Story) error {
	_, err := db.Exec(`
		INSERT INTO stories (user_id, content, image_url, bg_color)
		VALUES (?, ?, ?, ?)`,
		story.UserID, story.Content, story.ImageURL, story.BgColor,
	)
	return err
}

func GetActiveStories(db *sql.DB, authUserID string) ([]Storyapi, error) {
	query := `
          SELECT 
			s.id, 
			s.user_id,
			s.content, 
			s.image_url, 
			s.bg_color,
			s.created_at,
			s.expires_at,
            u.first_name, u.last_name, 
			u.image AS profile_image
        FROM stories s
        JOIN users u ON s.user_id = u.id
        WHERE 
			(s.expires_at IS NULL 
			OR DATETIME(s.expires_at) > CURRENT_TIMESTAMP)
			AND 
			(u.privacy = 'public'
			OR (u.privacy = 'private' AND EXISTS (
				SELECT 1 FROM followers f 
				WHERE f.user_id = s.user_id          
				AND f.follower_id = ?
			))
				OR (s.user_id = ?)
		)
        ORDER BY s.created_at ASC
    `
	rows, err := db.Query(query, authUserID, authUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stories []Storyapi
	for rows.Next() {
		var s Storyapi
		var content, imageURL, bg, firstName, lastName, profile sql.NullString
		var createdAt, expiresAt sql.NullString

		err := rows.Scan(&s.ID, &s.UserID, &content, &imageURL, &bg, &createdAt, &expiresAt,
			&firstName, &lastName, &profile)
		if err != nil {
			continue
		}

		s.Content = content.String
		s.ImageURL = imageURL.String
		s.BGColor = bg.String
		s.FirstName = firstName.String
		s.LastName = lastName.String
		s.Profile = profile.String

		if createdAt.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", createdAt.String)
			s.CreatedAt = t
		}

		if expiresAt.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", expiresAt.String)
			s.ExpiresAt = sql.NullTime{Time: t, Valid: true}
		}

		stories = append(stories, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stories, nil
}
