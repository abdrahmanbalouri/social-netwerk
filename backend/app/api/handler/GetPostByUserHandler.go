package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"social-network/app/helper"
	"social-network/app/repository/post"
	"social-network/pkg/db/sqlite"
)

func GetPostByUserHandler(w http.ResponseWriter, r *http.Request) {
if r.Method !=  http.MethodGet {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, " method not allowed ")
		return
}

	user, err1 := helper.AuthenticateUser(r)
	if err1 != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	offsetStr := parts[3]
	userId := r.URL.Query().Get("userId")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	posts, err := GetAllPostsByuser(userId, r, offset, user)
	//
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve posts")
		return
	}
	if posts == nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve posts")
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, posts)
}

func GetAllPostsByuser(authUserID string, r *http.Request, ofseet int, userId string) ([]map[string]interface{}, error) {
	var rows *sql.Rows
	var err error

	// userId, err := helper.AuthenticateUser(r)
	// if err != nil {
	// 	return nil, fmt.Errorf("authentication error: %v", err)
	// }
	if authUserID == "0" {
		authUserID = userId
	}
	limit := 10

	if authUserID != "" {
		rows, err = sqlite.Db.Query(`
			SELECT 
			  p.id, 
            p.user_id, 
            p.title, 
            p.content, 
            p.image_path,
            p.visibility,
            p.canseperivite,
            p.created_at, 
           u.first_name , u.last_name,
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
			GROUP BY p.id, p.user_id, p.title, p.content, p.image_path, p.created_at,  u.first_name , u.last_name, u.image
			ORDER BY p.created_at DESC
			LIMIT ? OFFSET ?;
		`, authUserID, authUserID, limit, ofseet)
	} else {
		rows, err = sqlite.Db.Query(`
	SELECT 
    p.id, 
    p.user_id, 
    p.title, 
    p.content, 
    p.image_path,
    p.visibility,
    p.canseperivite,
    p.created_at, 
    u.first_name , u.last_name,
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
GROUP BY p.id, p.user_id, p.title, p.content, p.image_path, p.created_at,  u.first_name , u.last_name, u.image
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
			first_name    string
			last_name     string

			profile       sql.NullString
			likeCount     int
			likedByUser   int
			commentsCount int
		)

		err := rows.Scan(&id, &userID, &title, &content, &imagePath, &visibility, &canseperivite, &createdAt, &first_name, &last_name, &privacy, &profile, &likeCount, &likedByUser, &commentsCount)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		post := map[string]interface{}{
			"id":            id,
			"user_id":       userID,
			"title":         title,
			"content":       content,
			"image_path":    post.NilIfEmpty(imagePath),
			"visibility":    visibility,
			"canseperivite": canseperivite,
			"privacy":       privacy,
			"created_at":    createdAt,
			"first_name":    first_name,
			"last_name":     last_name,

			"profile":        post.NilIfEmpty(profile),
			"like":           likeCount,
			"liked_by_user":  likedByUser > 0,
			"comments_count": commentsCount,
		}
		posts = append(posts, post)
	}
	if len(posts) == 0 {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return posts, nil
}
