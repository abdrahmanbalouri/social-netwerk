package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"social-network/internal/repository"
)

func GetStories(w http.ResponseWriter, r *http.Request) {
	// Query with JOIN and expiration filter directly in SQL (SQLite DATETIME handling)
	rows, err := repository.Db.Query(`
        SELECT s.id, s.user_id, s.content, s.image_url, s.bg_color, s.created_at, s.expires_at,
               u.nickname, u.image AS profile_image
        FROM stories s
        JOIN users u ON s.user_id = u.id
        WHERE s.expires_at IS NULL OR DATETIME(s.expires_at) > CURRENT_TIMESTAMP
        ORDER BY s.created_at DESC
    `)
	if err != nil {
		http.Error(w, "Failed to get stories", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stories []map[string]interface{}
	for rows.Next() {
		var (
			id, userID             string
			content, imageURL, bg  sql.NullString
			createdAt, expiresAt   sql.NullString
			nickname, profileImage sql.NullString
		)

		err := rows.Scan(&id, &userID, &content, &imageURL, &bg, &createdAt, &expiresAt, &nickname, &profileImage)
		if err != nil {
			// Log error if needed, but continue to next row
			fmt.Println("Scan error:", err)
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
			"nickname":   getNullString(nickname),
			"profile":    getNullString(profileImage),
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
		http.Error(w, "Error iterating stories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stories)
}

// Helper to handle sql.NullString safely (returns "" if invalid)
func getNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}