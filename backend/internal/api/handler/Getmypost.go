package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func Getmypost(w http.ResponseWriter, r *http.Request) {
	authUserID, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	userId := parts[3]
	offsetStr := parts[4]
	if userId == "0" {
		userId = authUserID
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
	limit := 10
	fmt.Println("2")
	var rows *sql.Rows

	if authUserID == userId {
		rows, err = repository.Db.Query(`
			SELECT 
				p.id, 
				p.user_id, 
				p.title, 
				p.content, 
				p.image_path,
				p.visibility,
				p.canseperivite,
				p.created_at, 
				u.nickname,
				u.privacy,
				u.image AS profile,
				COUNT(DISTINCT l.id) AS like_count,
				COUNT(DISTINCT CASE WHEN l.user_id = ? THEN l.id END) AS liked_by_user,
				COUNT(DISTINCT c.id) AS comments_count
			FROM posts p
			JOIN users u ON p.user_id = u.id
			LEFT JOIN likes l ON p.id = l.liked_item_id AND l.liked_item_type = 'post'
			LEFT JOIN comments c ON p.id = c.post_id
			WHERE p.user_id = ?
			GROUP BY 
				p.id, p.user_id, p.title, p.content, p.image_path, p.visibility,
				p.canseperivite, p.created_at, u.nickname, u.privacy, u.image
			ORDER BY p.created_at DESC
			LIMIT ? OFFSET ?;
		`, authUserID, userId, limit, offset)
	} else {
		rows, err = repository.Db.Query(`
			SELECT 
				p.id, 
				p.user_id, 
				p.title, 
				p.content, 
				p.image_path,
				p.visibility,
				p.canseperivite,
				p.created_at, 
				u.nickname,
				u.privacy,
				u.image AS profile,
				COUNT(DISTINCT l.id) AS like_count,
				COUNT(DISTINCT CASE WHEN l.user_id = ? THEN l.id END) AS liked_by_user,
				COUNT(DISTINCT c.id) AS comments_count
			FROM posts p
			JOIN users u ON p.user_id = u.id
			LEFT JOIN likes l ON p.id = l.liked_item_id AND l.liked_item_type = 'post'
			LEFT JOIN comments c ON p.id = c.post_id
			WHERE p.user_id = ?
			  AND (
					u.privacy = 'public'
					OR (
						u.privacy = 'private'
						AND EXISTS (
							SELECT 1 FROM followers f
							WHERE f.user_id = p.user_id
							  AND f.follower_id = ?
						)
					)
				)
			GROUP BY 
				p.id, p.user_id, p.title, p.content, p.image_path, p.visibility,
				p.canseperivite, p.created_at, u.nickname, u.privacy, u.image
			ORDER BY p.created_at DESC
			LIMIT ? OFFSET ?;
		`, authUserID, userId, authUserID, limit, offset)
	}

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Query error: "+err.Error())
		return
	}
	defer rows.Close()

	var posts []map[string]interface{}

	for rows.Next() {
		var (
			id, userID, title, content, visibility, canseperivite, createdAt, privacy, nickname string
			imagePath, profile                                                                  sql.NullString
			likeCount, likedByUser, commentsCount                                               int
		)

		err := rows.Scan(
			&id, &userID, &title, &content, &imagePath,
			&visibility, &canseperivite, &createdAt, &nickname, &privacy,
			&profile, &likeCount, &likedByUser, &commentsCount,
		)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Scan error: "+err.Error())
			return
		}

		post := map[string]interface{}{
			"id":             id,
			"user_id":        userID,
			"title":          title,
			"content":        content,
			"image_path":     nilIfEmpty(imagePath),
			"visibility":     visibility,
			"canseperivite":  canseperivite,
			"privacy":        privacy,
			"created_at":     createdAt,
			"author":         nickname,
			"profile":        nilIfEmpty(profile),
			"like":           likeCount,
			"liked_by_user":  likedByUser > 0,
			"comments_count": commentsCount,
		}
		posts = append(posts, post)
	}

	if len(posts) == 0 {
		helper.RespondWithError(w, http.StatusNotFound, "No posts found")
		return
	}

	if err := rows.Err(); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Rows error: "+err.Error())
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, posts)
}

func nilIfEmpty(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}
