CREATE OR REPLACE FUNCTION cleanup_notifications(notif notifications)
    RETURNS VOID AS
$$
    DECLARE
        mark_for_delete BOOLEAN;    
    BEGIN
        IF notif.subject_type != '' THEN
            EXECUTE 'SELECT CASE WHEN id IS NULL THEN true ELSE false END FROM ' || notif.subject_type || ' WHERE id = ' || notif.subject_id INTO mark_for_delete;
        END IF;

        IF notif.object_type != '' AND NOT COALESCE(mark_for_delete, true) THEN
            EXECUTE 'SELECT CASE WHEN id IS NULL THEN true ELSE false END FROM ' || notif.object_type || ' WHERE id = ' || notif.object_id INTO mark_for_delete;
        END IF;

        IF notif.indirect_object_type != '' AND NOT COALESCE(mark_for_delete, true) THEN
            EXECUTE 'SELECT CASE WHEN id IS NULL THEN true ELSE false END FROM ' || notif.indirect_object_type || ' WHERE id = ' || notif.indirect_object_id INTO mark_for_delete;
        END IF;

        IF COALESCE(mark_for_delete, true) THEN
            DELETE FROM notifications
            WHERE id = notif.id;
        END IF;
    END;
$$ LANGUAGE plpgsql;

DO $$
DECLARE
BEGIN
    PERFORM cleanup_notifications(n.*)
    FROM notifications AS n;
END $$ LANGUAGE plpgsql;

DROP FUNCTION IF EXISTS cleanup_notifications;
