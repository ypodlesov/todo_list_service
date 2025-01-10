CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL,
    title VARCHAR(128),
    description VARCHAR(4096),
    status SMALLINT, -- 1 (opened), 2 (closed)
    priority INTEGER,
    user_id INTEGER,
    creation_ts TIMESTAMP DEFAULT 'now'
);

CREATE INDEX IF NOT EXISTS tasks_user_id_idx ON tasks (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS tasks_id_idx ON tasks (id);