package service

import "social-network/internal/repository"

func LogoutUser(token string) error {
	return repository.DeleteSession(repository.Db , token)
}
