package service

import (
	"social-network/app/repository/model"
	"social-network/pkg/db/sqlite"
)

// FetchUserGallery calls the database layer
func FetchUserGallery(userID string) ([]model.PostGallery, error) {
	return model.GetUserGallery(sqlite.Db, userID)
}
