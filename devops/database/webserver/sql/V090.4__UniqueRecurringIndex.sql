DROP INDEX IF EXISTS recurring_tasks_event_id_idx;
CREATE UNIQUE INDEX ON recurring_tasks(event_id);
