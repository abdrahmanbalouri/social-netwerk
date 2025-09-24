
-- +migrate Up
CREATE TABLE IF NOT EXISTS group_members (
    user_id TEXT NOT NULL,
    group_id TEXT NOT NULL,
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);
-- +migrate Down
DROP TABLE IF EXISTS group_members;