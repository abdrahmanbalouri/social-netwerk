package handlers

import (
	"database/sql"
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func GetVedioHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	rows, err := repository.Db.Query(`
		SELECT 
			p.id, 
			p.user_id, 
			p.title, 
			p.content, 
			p.image_path,
			p.visibility,
			p.canseperivite,
			p.created_at, 
			u.first_name, u.last_name,
			u.privacy,
			u.image AS profile,
			COALESCE(COUNT(DISTINCT l.id), 0) AS like_count,
			COALESCE(COUNT(DISTINCT CASE WHEN l.user_id = ? THEN l.id END), 0) AS liked_by_user,
			COALESCE(COUNT(DISTINCT c.id), 0) AS comments_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN likes l ON p.id = l.liked_item_id AND l.liked_item_type = 'post'
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE 
			p.image_path IS NOT NULL
			AND (
				LOWER(p.image_path) LIKE '%.mp4' OR
				LOWER(p.image_path) LIKE '%.webm' OR
				LOWER(p.image_path) LIKE '%.ogg' OR
				LOWER(p.image_path) LIKE '%.mov'
			)
			AND (
				p.user_id = ? 
				OR (p.visibility = 'public' AND u.privacy = 'public')
				OR (p.visibility = 'public' AND u.privacy = 'private' AND EXISTS (
					SELECT 1 FROM followers f 
					WHERE f.user_id = p.user_id AND f.follower_id = ?
				))
				OR (p.visibility = 'almost_private' AND EXISTS (
					SELECT 1 FROM followers f 
					WHERE f.user_id = p.user_id AND f.follower_id = ?
				))
				OR (p.visibility = 'private' AND EXISTS (
					SELECT 1 FROM allowed_followers af 
					WHERE af.allowed_user_id = ? AND af.user_id = p.user_id
				))
			)
		GROUP BY p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, 
				 u.first_name, u.last_name, u.image
		ORDER BY p.created_at DESC
	`, userId, userId, userId, userId, userId)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var (
			id            string
			userID        string
			title         string
			content       string
			imagePath     sql.NullString
			visibility    string
			canseperivite string
			createdAt     string
			first_name    string
			last_name     string
			privacy       string
			profile       sql.NullString
			likeCount     int
			likedByUser   int
			commentsCount int
		)

		err := rows.Scan(
			&id, &userID, &title, &content, &imagePath,
			&visibility, &canseperivite, &createdAt,
			&first_name, &last_name, &privacy, &profile,
			&likeCount, &likedByUser, &commentsCount,
		)
		if err != nil {
			continue
		}

		if !imagePath.Valid {
			continue
		}

		post := map[string]interface{}{
			"id":             id,
			"user_id":        userID,
			"title":          title,
			"content":        content,
			"image_path":     NilIfEmpty(imagePath),
			"visibility":     visibility,
			"canseperivite":  canseperivite,
			"privacy":        privacy,
			"created_at":     createdAt,
			"first_name":     first_name,
			"last_name":      last_name,
			"profile":        NilIfEmpty(profile),
			"like":           likeCount,
			"liked_by_user":  likedByUser > 0,
			"comments_count": commentsCount,
		}
		posts = append(posts, post)
	}

	helper.RespondWithJSON(w, http.StatusOK, posts)
}

func NilIfEmpty(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}
