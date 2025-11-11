package service

import (
	"social-network/app/repository/model"
	"social-network/pkg/db/sqlite"
)

func LogoutUser(token string) error {
	return model.DeleteSession(sqlite.Db, token)
}
