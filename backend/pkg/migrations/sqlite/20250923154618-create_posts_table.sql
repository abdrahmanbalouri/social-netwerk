
-- +migrate Up
CREATE TABLE IF NOT EXISTS posts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    visibility TEXT DEFAULT 'public',
    image_path TEXT,
    canseperivite  TEXT DEFAULT "",
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +migrate Down
DROP TABLE IF EXISTS posts;