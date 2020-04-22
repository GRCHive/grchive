ALTER TABLE scheduled_tasks
DROP COLUMN task_date;

ALTER TABLE scheduled_tasks
ADD COLUMN task_data JSONB NOT NULL;
