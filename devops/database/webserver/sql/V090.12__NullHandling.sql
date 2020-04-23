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

        IF one_time.event_id IS NULL THEN
            one_time := NULL;
        END IF;

        SELECT * INTO recurring
        FROM recurring_tasks
        WHERE event_id = NEW.id;

        IF recurring.event_id IS NULL THEN
            recurring := NULL;
        END IF;

        PERFORM pg_notify(
            'scheduledtaskchange',
            jsonb_build_object(
                'Task', NEW,
                'OneTime', one_time,
                'Recurring', recurring 
            ) #>> '{}'
        );
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
