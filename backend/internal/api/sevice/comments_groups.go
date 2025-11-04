package service

import (
	"errors"
	"time"
	"social-network/internal/helper"
	"social-network/internal/repository"
)

func GetCommentsGroup(userID, groupID, postID string, offset int) ([]map[string]interface{}, error) {
	if err := helper.CheckUserInGroup(userID, groupID); err != nil {
		return nil, errors.New("you are not a member of this group")
	}

	type Comment struct {
		ID        string
		Content   string
		CreatedAt time.Time
		FirstName string
		LastName  string
		MediaPath string
	}

	repoComments, err := repository.FetchCommentsGroup( repository.Db ,postID, offset)
	if err != nil {
		return nil, err
	}

	// تحويل repo struct إلى map
	var comments []map[string]interface{}
	for _, c := range repoComments {
		comments = append(comments, map[string]interface{}{
			"id":         c.ID,
			"content":    c.Content,
			"created_at": c.CreatedAt.Format(time.RFC3339),
			"first_name": c.FirstName,
			"last_name":  c.LastName,
			"media_path": c.MediaPath,
		})
	}

	return comments, nil
}
