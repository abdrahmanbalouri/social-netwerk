package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	

	//"social-network/internal/database"
	"social-network/internal/helper"
	"social-network/internal/repository"
)

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 { // [0]= "", [1]=api, [2]=Getcomments, [3]=postID, [4]=offset
		helper.RespondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	postID := parts[3]
	offsetStr := parts[4]

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
	_, err = helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// can, err := helper.CanViewComments(userID, postID)
	// if err != nil {
	// 	helper.RespondWithError(w, http.StatusInternalServerError, "Error checking permissions")
	// 	return
	// }
	// if !can {
	// 	helper.RespondWithError(w, http.StatusForbidden, "You do not have permission to view comments on this post")
	// 	return
	// }

	rows, err := repository.Db.Query(`
        SELECT c.id, c.content, c.created_at, u.nickname
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = ?
        ORDER BY c.created_at desc
        LIMIT 10 OFFSET ?`, postID, offset)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch comments")
		return
	}
	defer rows.Close()

	type Comment struct {
		ID        string    `json:"id"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
		Author    string    `json:"author"`
	}

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.Author)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Failed to process comments")
			return
		}
		comments = append(comments, comment)
	}

	helper.RespondWithJSON(w, http.StatusOK, comments)
}
