-- +migrate Up

CREATE TABLE eventsAction_new (
 id INTEGER PRIMARY KEY AUTOINCREMENT,    user_id TEXT NOT NULL,
    event_id TEXT NOT NULL,
    action  TEXT DEFAULT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
);



.DROP TABLE event_Actions;

ALTER TABLE eventsAction_new RENAME TO event_Actions;
-- +migrate Down
DROP TABLE IF EXISTS event_Actions;
