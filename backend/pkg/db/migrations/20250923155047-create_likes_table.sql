
-- +migrate Up
CREATE TABLE IF NOT EXISTS likes (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    liked_item_id TEXT NOT NULL,
    liked_item_type TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
    FOREIGN KEY (liked_item_id) REFERENCES posts(id)
);

-- +migrate Down
DROP TABLE IF EXISTS likes;