
-- +migrate Up
CREATE TABLE followers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id  TEXT ,
    follower_id TEXT ,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, follower_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(follower_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE followers;
