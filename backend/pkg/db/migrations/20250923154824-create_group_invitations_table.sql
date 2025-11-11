
-- +migrate Up
CREATE TABLE IF NOT EXISTS group_invitations (
    id TEXT PRIMARY KEY,
    group_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    invited_by_user_id TEXT,
    request_type TEXT NOT NULL CHECK (request_type IN ('invitation', 'join')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (invited_by_user_id) REFERENCES users(id),
    CHECK (
        (request_type = 'invitation' AND invited_by_user_id IS NOT NULL) OR 
        (request_type = 'join' AND invited_by_user_id IS NULL)
    )
);
-- +migrate Down
DROP TABLE IF EXISTS group_invitations;