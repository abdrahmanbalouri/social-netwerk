package post

import (
	"database/sql"
	"fmt"
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func GetAllPosts(authUserID string, r *http.Request, ofseet int) ([]map[string]interface{}, error) {
	var rows *sql.Rows
	var err error

	userId, err := helper.AuthenticateUser(r)
	if err != nil {
		return nil, fmt.Errorf("authentication error: %v", err)
	}
	if authUserID == "0" {
		authUserID = userId
	}
	limit := 10

	if authUserID != "" {
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
			GROUP BY p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, u.nickname, u.image
			ORDER BY p.created_at DESC
			LIMIT ? OFFSET ?;
		`, authUserID, authUserID, limit, ofseet)
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
WHERE 
    p.user_id = ? 
    OR (p.visibility = 'public' AND u.privacy = 'public')
    OR (p.visibility = 'public' AND u.privacy = 'private' AND EXISTS (
        SELECT 1 FROM followers f 
        WHERE f.user_id = ?         
          AND f.follower_id = p.user_id  
    ))
    OR (p.visibility = 'almost_private' AND EXISTS (
        SELECT 1 FROM followers f 
        WHERE f.user_id = ? 
          AND f.follower_id = p.user_id
    ))
    OR (p.visibility = 'private' AND EXISTS (
        SELECT 1 FROM allowed_followers af 
        WHERE af.allowed_user_id = ? 
          AND af.user_id = p.user_id
    ))
GROUP BY p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, u.nickname, u.image
ORDER BY p.created_at DESC
LIMIT ? OFFSET ?;
`, userId, userId, userId, userId, userId, limit, ofseet)
	}

	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
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
			privacy       string
			nickname      string
			profile       sql.NullString
			likeCount     int
			likedByUser   int
			commentsCount int
		)

		err := rows.Scan(&id, &userID, &title, &content, &imagePath, &visibility, &canseperivite, &createdAt, &nickname, &privacy, &profile, &likeCount, &likedByUser, &commentsCount)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
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
	if len(posts) == 0{
		return  nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return posts, nil
}

// Helper to convert sql.NullString to interface{}
func nilIfEmpty(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}