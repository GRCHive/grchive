CREATE OR REPLACE FUNCTION trigger_scheduled_task_change_notification()
    RETURNS trigger AS
$$
    DECLARE
        one_time one_time_tasks;
        recurring recurring_tasks;
    BEGIN
        SELECT * INTO one_time
        FROM one_time_tasks
        WHERE event_id = NEW.id;

        SELECT * INTO recurring
        FROM recurring_tasks
        WHERE event_id = NEW.id;

        PERFORM pg_notify(
            'scheduledtaskchange',
            jsonb_build_object(
                'Task', NEW,
                'OneTime', one_time,
                'Recurring', recurring 
            ) #>> '{}'
        );
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS scheduled_task_change_notification ON scheduled_tasks;
CREATE CONSTRAINT TRIGGER scheduled_task_change_notification
    AFTER INSERT OR UPDATE ON scheduled_tasks
    DEFERRABLE INITIALLY DEFERRED
    FOR EACH ROW
    EXECUTE FUNCTION trigger_scheduled_task_change_notification();

---

CREATE OR REPLACE FUNCTION trigger_scheduled_task_delete_notification()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM pg_notify(
            'scheduledtaskdelete',
            jsonb_build_object(
                'Task', OLD
            ) #>> '{}'
        );
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS scheduled_task_delete_notification ON scheduled_tasks;
CREATE TRIGGER scheduled_task_delete_notification
    BEFORE DELETE ON scheduled_tasks
    FOR EACH ROW
    EXECUTE FUNCTION trigger_scheduled_task_delete_notification();
