CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(1024),
    status SMALLINT,
    user_id INTEGER,
    creation_ts TIMESTAMP DEFAULT 'now'
);

CREATE INDEX IF NOT EXISTS user_id_idx ON tasks (user_id);