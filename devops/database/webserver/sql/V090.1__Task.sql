CREATE TABLE supported_task_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256)
);

INSERT INTO supported_task_types(id, name)
VALUES
    (1, 'RabbitMQ');

ALTER TABLE scheduled_tasks
ADD COLUMN task_type INTEGER NOT NULL REFERENCES supported_task_types(id);

ALTER TABLE scheduled_tasks
ADD COLUMN task_date JSONB NOT NULL;
