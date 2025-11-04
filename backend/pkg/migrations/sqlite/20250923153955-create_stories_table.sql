-- +migrate Up
CREATE TABLE IF NOT EXISTS stories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    content TEXT,
    image_url TEXT,
    bg_color TEXT DEFAULT '#000000',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME DEFAULT (DATETIME('now', '+3 minutes')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +migrate Down
DROP TABLE IF EXISTS stories;