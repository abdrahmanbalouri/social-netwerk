package service

import (
	"social-network/internal/repository"
	"social-network/internal/repository/model"
)

func LogoutUser(token string) error {
	return model.DeleteSession(repository.Db , token)
}
