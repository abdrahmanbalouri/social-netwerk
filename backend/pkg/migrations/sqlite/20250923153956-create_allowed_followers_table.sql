-- +migrate Up
CREATE TABLE IF NOT EXISTS allowed_followers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,          
    post_id TEXT NOT NULL,          
    allowed_user_id TEXT NOT NULL,  
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (allowed_user_id) REFERENCES users(id)
);

-- +migrate Down
DROP TABLE IF EXISTS allowed_followers;
