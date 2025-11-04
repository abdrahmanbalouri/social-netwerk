package service

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func LikesGroupService(r *http.Request, userID string) (map[string]string, error) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		return nil, fmt.Errorf("post ID and group ID are required")
	}
	postID := pathParts[3]
	groupID := pathParts[4]

	// Check if user is in the group
	if err := helper.CheckUserInGroup(userID, groupID); err != nil {
		return nil, fmt.Errorf("you are not a member of this group")
	}

	exists, err := repository.DoesPostExistInGroup(postID, groupID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	if !exists {
		return nil, fmt.Errorf("post not found in this group")
	}

	existingLikeID, err := repository.GetExistingLikeGroup(userID, postID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("database error")
	}

	if existingLikeID != "" {
		// Remove like
		if err := repository.RemoveLikeGroup(existingLikeID); err != nil {
			return nil, fmt.Errorf("failed to remove like")
		}
		return map[string]string{"message": "Like removed"}, nil
	}

	// Add like
	likeId := helper.GenerateUUID()
	if err := repository.AddLikeGroup(likeId.String(), userID, postID, time.Now()); err != nil {
		return nil, fmt.Errorf("failed to add like")
	}

	return map[string]string{"message": "Like added"}, nil
}
