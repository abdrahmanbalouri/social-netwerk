-- +migrate Up
CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id TEXT NOT NULL,
    receiver_id TEXT NOT NULL,
    type TEXT NOT NULL,             
    message TEXT DEFAULT "",
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);

-- +migrate Down
DROP TABLE IF EXISTS notifications;
