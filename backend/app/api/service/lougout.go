package service

import (
	"social-network/app/repository"
	"social-network/app/repository/model"
)

func LogoutUser(token string) error {
	return model.DeleteSession(repository.Db , token)
}
