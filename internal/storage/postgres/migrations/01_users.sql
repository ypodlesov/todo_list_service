CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(64),
    password VARCHAR(256),
    email VARCHAR(128)
);

CREATE UNIQUE INDEX IF NOT EXISTS username_idx ON users (username);