
-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    nickname TEXT DEFAULT "",
    date_birth INTEGER NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    image TEXT DEFAULT "default.png",
    cover TEXT DEFAULT "",
    about TEXT NOT NULL,
    privacy TEXT NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);

-- +migrate Down
DROP TABLE IF EXISTS users;