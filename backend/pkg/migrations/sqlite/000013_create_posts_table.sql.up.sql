CREATE TABLE
    posts (
        id TEXT PRIMARY KEY,
        user_id INT NOT NULL,
        group_id INT DEFAULT NULL,
        content TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (id),
        FOREIGN KEY (group_id) REFERENCES groups (id)
    );