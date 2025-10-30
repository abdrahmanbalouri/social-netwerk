-- +migrate Up
CREATE TABLE IF NOT EXISTS events (
    id TEXT PRIMARY KEY,
    group_id TEXT NOT NULL,
    title TEXT DEFAULT NULL,
    description  TEXT NOT NULL,
    time DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);
-- +migrate Down
DROP TABLE IF EXISTS events;