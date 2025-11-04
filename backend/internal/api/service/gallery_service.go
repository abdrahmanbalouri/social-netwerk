package service

import (
	"social-network/internal/repository"
	"social-network/internal/repository/model"
)

// FetchUserGallery calls the database layer
func FetchUserGallery(userID string) ([]model.PostGallery, error) {
	return model.GetUserGallery(repository.Db, userID)
}
