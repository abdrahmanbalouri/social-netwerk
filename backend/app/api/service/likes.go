package service

import (
	"database/sql"
	"time"

	"social-network/app/helper"
	"social-network/app/repository/model"
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
		ID:            helper.GenerateUUID().String(),
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
