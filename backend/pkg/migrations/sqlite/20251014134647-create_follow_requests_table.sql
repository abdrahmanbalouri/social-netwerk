
-- +migrate Up
CREATE TABLE follow_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    follower_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, follower_id)
);

-- +migrate Down
DROP TABLE follow_requests;