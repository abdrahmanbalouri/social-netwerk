package service

import (
	"social-network/app/repository"
	"social-network/app/repository/model"
)

// FetchUserGallery calls the database layer
func FetchUserGallery(userID string) ([]model.PostGallery, error) {
	return model.GetUserGallery(repository.Db, userID)
}
