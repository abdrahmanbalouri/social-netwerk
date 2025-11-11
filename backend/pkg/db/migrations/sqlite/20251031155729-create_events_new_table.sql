-- +migrate Up

CREATE TABLE events_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id TEXT NOT NULL,
    title TEXT DEFAULT NULL,
    description TEXT NOT NULL,
    time DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

INSERT INTO events_new (group_id, title, description, time, created_at)
SELECT group_id, title, description, time, created_at FROM events;

DROP TABLE events;

ALTER TABLE events_new RENAME TO events;
-- +migrate Down
DROP TABLE IF EXISTS events;
