package repository

import "database/sql"

func SearchUsersInDB(db *sql.DB, search string) ([]map[string]any, error) {
    rows, err := db.Query(`
        SELECT id, nickname, first_name, last_name, image 
        FROM users 
        WHERE lower(nickname) LIKE ? OR lower(first_name) LIKE ? OR lower(last_name) LIKE ?
        LIMIT 10
    `, search, search, search)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []map[string]any
    for rows.Next() {
        var id, nickname, firstName, lastName, image string
        if err := rows.Scan(&id, &nickname, &firstName, &lastName, &image); err != nil {
            continue
        }
        users = append(users, map[string]any{
            "id":         id,
            "nickname":   nickname,
            "first_name": firstName,
            "last_name":  lastName,
            "image":      image,
        })
    }
    return users, nil
}
