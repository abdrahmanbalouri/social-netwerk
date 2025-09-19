-- +migrate Up
CREATE TABLE IF NOT EXISTS comments (
 	id TEXT PRIMARY KEY,
	post_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
	content TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
	FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE

);

-- +migrate Down
DROP TABLE IF EXISTS comments;
