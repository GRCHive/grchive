ALTER TABLE scheduled_tasks
ADD COLUMN user_id_new BIGINT REFERENCES users(id) ON DELETE CASCADE;

UPDATE scheduled_tasks
SET user_id_new = user_id;

ALTER TABLE scheduled_tasks
ALTER COLUMN user_id_new SET NOT NULL;

ALTER TABLE scheduled_tasks
DROP COLUMN user_id;
