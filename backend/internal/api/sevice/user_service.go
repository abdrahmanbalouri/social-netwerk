package service

import (
	"database/sql"
	"strings"

	"social-network/internal/repository"
)

func SearchUsers(query string) ([]map[string]any, error) {
	if query == "" {
		return []map[string]any{}, nil
	}
	search := "%" + strings.ToLower(query) + "%"
	return repository.SearchUsersInDB(repository.Db, search)
}

func GetUsers(db *sql.DB, currentUserID string) ([]repository.User, error) {
	return repository.FetchAllUsers(db, currentUserID)
}
