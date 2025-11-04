package repository

import (
	"time"

	"social-network/internal/utils"
)

func InsertGroupPost(id, userID, groupID, title, content, imagePath string, createdAt time.Time) error {
	_, err := Db.Exec(`
		INSERT INTO group_posts (id, user_id, group_id, title, content, image_path, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, userID, groupID, title, content, imagePath, createdAt,
	)
	return err
}

func GetGroupPostByID(postID, userID string) (map[string]interface{}, error) {
	post := make(map[string]interface{})
	query := `
	SELECT  
		gp.id, gp.user_id, gp.title, gp.content, gp.image_path, gp.created_at, 
		u.first_name, u.last_name, u.image AS profile, 
		COUNT(DISTINCT l.id) AS like_count, 
		COUNT(DISTINCT CASE WHEN l.user_id = ? THEN l.id END) AS liked_by_user, 
		COUNT(DISTINCT c.id) AS comments_count 
	FROM group_posts gp 
	JOIN users u ON gp.user_id = u.id 
	LEFT JOIN likesgroups l ON gp.id = l.liked_item_id AND l.liked_item_type = 'post' 
	LEFT JOIN comments_groups c ON gp.id = c.post_id 
	WHERE gp.id = ?
	GROUP BY gp.id, gp.user_id, gp.title, gp.content, gp.image_path, gp.created_at, u.first_name, u.last_name, u.image
	LIMIT 1;
	`

	var (
		id, userId, title, content, imagePath, firstName, lastName, profile string
		createdAt                                                           time.Time
		likeCount, likedByUser, commentsCount                               int
	)

	err := Db.QueryRow(query, userID, postID).Scan(
		&id, &userId, &title, &content, &imagePath, &createdAt,
		&firstName, &lastName, &profile, &likeCount, &likedByUser, &commentsCount,
	)
	if err != nil {
		return nil, err
	}

	post["id"] = id
	post["user_id"] = userId
	post["title"] = title
	post["content"] = content
	post["image_path"] = imagePath
	post["created_at"] = createdAt
	post["first_name"] = firstName
	post["last_name"] = lastName
	post["profile"] = profile
	post["like_count"] = likeCount
	post["liked_by_user"] = likedByUser > 0
	post["comments_count"] = commentsCount

	return post, nil
}

func GetAllGroupPosts(groupID, userID string) ([]utils.GroupPost, error) {
	query := `
	SELECT 
		gp.id, 
		gp.user_id, 
		gp.title, 
		gp.content, 
		gp.image_path, 
		gp.created_at, 
		u.first_name, u.last_name, u.image AS profile,
		COUNT(DISTINCT l.id) AS like_count,
		COUNT(DISTINCT CASE WHEN l.user_id = ? THEN l.id END) AS liked_by_user,
		COUNT(DISTINCT c.id) AS comments_count
	FROM group_posts gp
	JOIN users u ON gp.user_id = u.id
	LEFT JOIN likesgroups l ON gp.id = l.liked_item_id AND l.liked_item_type = 'post'
	LEFT JOIN comments_groups c ON gp.id = c.post_id
	WHERE gp.group_id = ?
	GROUP BY 
		gp.id, gp.user_id, gp.title, gp.content, gp.image_path, gp.created_at, 
		u.first_name, u.last_name, u.image
	ORDER BY gp.created_at DESC;`

	rows, err := Db.Query(query, userID, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []utils.GroupPost
	for rows.Next() {
		var p utils.GroupPost
		err := rows.Scan(
			&p.ID, &p.UserID, &p.Title, &p.Content, &p.ImagePath,
			&p.CreatedAt, &p.FirstName, &p.LastName, &p.Profile,
			&p.LikeCount, &p.LikedByUser, &p.CommentsCount,
		)
		if err != nil {
			return nil, err
		}
		p.LikedByUser = p.LikeCount > 0
		posts = append(posts, p)
	}
	return posts, nil
}
