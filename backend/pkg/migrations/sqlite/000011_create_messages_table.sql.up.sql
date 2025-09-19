CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    sender_id INTEGER NOT NULL,
    content TEXT,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    group_id INTEGER DEFAULT NULL,
    receiver_id INTEGER DEFAULT NULL,
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (receiver_id) REFERENCES users(id)
);