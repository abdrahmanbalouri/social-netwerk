package service

import (
	"database/sql"
	"strings"

	"social-network/app/repository/model"
	"social-network/pkg/db/sqlite"
)

func SearchUsers(query string) ([]map[string]any, error) {
	if query == "" {
		return []map[string]any{}, nil
	}
	search := "%" + strings.ToLower(query) + "%"
	return model.SearchUsersInDB(sqlite.Db, search)
}

func GetUsers(db *sql.DB, currentUserID string) ([]model.User, error) {
	return model.FetchAllUsers(db, currentUserID)
}
