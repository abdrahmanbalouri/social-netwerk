CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nickname TEXT DEFAULT "",
    date_birth INTEGER NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    image TEXT NOT null ,
    about TEXT NOT NULL ,
    privacy TEXT NOT NULL ,
    created_at INTEGER NOT NULL
);