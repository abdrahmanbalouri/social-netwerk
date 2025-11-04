package service

import (
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
