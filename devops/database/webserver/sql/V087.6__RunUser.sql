DELETE FROM script_runs;

ALTER TABLE script_runs
ADD COLUMN user_id BIGINT NOT NULL REFERENCES users(id);
