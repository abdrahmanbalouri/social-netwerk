package service

import (
	"database/sql"
	"time"

	"social-network/internal/helper"
	"social-network/internal/repository/model"

	"github.com/google/uuid"
)

func TogglePostLike(db *sql.DB, userID, postID string) (string, error) {
	ok, err := helper.Canshowdata(userID, postID)
	if !ok {
		return "", err
	}
	exists, likeID, err := model.UserLikedPost(db, userID, postID)
	if err != nil {
		return "", err
	}

	if exists {
		// Unlike
		err = model.DeleteLike(db, likeID)
		if err != nil {
			return "", err
		}
		return "removed", nil
	}

	// Add like
	newLike := model.Like{
		ID:            uuid.New().String(),
		UserID:        userID,
		LikedItemID:   postID,
		LikedItemType: "post",
		CreatedAt:     time.Now(),
	}
	err = model.InsertLike(db, newLike)
	if err != nil {
		return "", err
	}
	return "added", nil
}
