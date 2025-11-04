-- +migrate Up
CREATE TABLE IF NOT EXISTS event_Actions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    event_id TEXT NOT NULL,
    action  TEXT DEFAULT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE

);
-- +migrate Down
DROP TABLE IF EXISTS event_Actions;