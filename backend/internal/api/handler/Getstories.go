package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func GetStories(w http.ResponseWriter, r *http.Request) {
	id, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	rows, err := repository.Db.Query(`
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
    `, id, id)
	fmt.Println(err)
	if err != nil {
		helper.RespondWithError(w, http.StatusAccepted, "not story yet")
		return
	}
	defer rows.Close()

	var stories []map[string]interface{}

	for rows.Next() {
		var (
			id, userID                          string
			content, imageURL, bg               sql.NullString
			createdAt, expiresAt                sql.NullString
			first_name, last_name, profileImage sql.NullString
		)

		err := rows.Scan(&id, &userID, &content, &imageURL, &bg, &createdAt, &expiresAt, &first_name, &last_name, &profileImage)
		if err != nil {

			continue
		}

		story := map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"content":    getNullString(content),
			"image_url":  getNullString(imageURL),
			"bg_color":   getNullString(bg),
			"created_at": getNullString(createdAt),
			"expires_at": getNullString(expiresAt),
			"first_name": getNullString(first_name),
			"last_name":  getNullString(last_name),

			"profile": getNullString(profileImage),
		}

		// Optional manual filter (for testing or if DB timezone issues)
		if expiresAt.Valid {
			expTime, parseErr := time.Parse("2006-01-02 15:04:05", expiresAt.String)
			if parseErr == nil && expTime.Before(time.Now()) {
				continue // skip expired
			}
		}

		stories = append(stories, story)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		helper.RespondWithError(w, http.StatusAccepted, "not story yet")
		return
	}
 
	
	helper.RespondWithJSON(w, http.StatusOK, stories)

}
// Helper to handle sql.NullString safely (returns "" if invalid)
func getNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
