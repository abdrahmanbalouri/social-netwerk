
-- +migrate Up
CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    sender_id TEXT NOT NULL,
    content TEXT,
    sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    group_id TEXT,
    receiver_id TEXT,
<<<<<<< HEAD
    image TEXT,
=======
>>>>>>> fix-groups-errors
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (receiver_id) REFERENCES users(id)
);
-- +migrate Down
DROP TABLE IF EXISTS messages;