package service

import "social-network/internal/repository"

// FetchUserGallery calls the database layer
func FetchUserGallery(userID string) ([]repository.PostGallery, error) {
	return repository.GetUserGallery(repository.Db, userID)
}
