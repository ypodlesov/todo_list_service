CREATE TABLE IF NOT EXISTS task_actions (
    id SERIAL PRIMARY KEY,
    action_type SMALLINT,
    user_id INTEGER,
    task_id INTEGER,
    ts TIMESTAMP DEFAULT 'now'
);

CREATE INDEX IF NOT EXISTS task_id_idx ON task_actions (task_id);