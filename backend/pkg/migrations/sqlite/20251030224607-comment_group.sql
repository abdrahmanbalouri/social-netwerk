
-- +migrate Up
CREATE TABLE IF NOT EXISTS comments_groups (
    id TEXT PRIMARY KEY,
    post_id  TEXT NOT NULL,
    user_id TEXT NOT NULL,
    content TEXT,
    media_path TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES group_posts(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS comments_groups;