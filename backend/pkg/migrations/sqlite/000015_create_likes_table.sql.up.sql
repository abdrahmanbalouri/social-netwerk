CREATE TABLE likes (
    id TEXT PRIMARY KEY,
    user_id INT NOT NULL,
    liked_item_id INT NOT NULL,
    liked_item_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);