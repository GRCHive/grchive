DROP FUNCTION IF EXISTS trigger_scheduled_task_change_notification CASCADE;
DROP TRIGGER IF EXISTS scheduled_task_change_notification ON scheduled_tasks;

DROP FUNCTION IF EXISTS trigger_scheduled_task_delete_notification CASCADE;
DROP TRIGGER IF EXISTS scheduled_task_delete_notification ON scheduled_tasks;

