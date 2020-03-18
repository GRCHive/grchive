CREATE OR REPLACE FUNCTION trigger_document_request_change_notification()
    RETURNS trigger AS
$$
    DECLARE
        action_user users;
    BEGIN
        SELECT u.* INTO action_user 
        FROM users AS u
        INNER JOIN postgres_oid_to_users AS lnk
            ON lnk.user_id = u.id
        INNER JOIN pg_roles AS rl
            ON rl.oid = lnk.pg_oid
        WHERE rolname = current_user;

        IF NEW.assignee != COALESCE(OLD.assignee, -1) THEN
            PERFORM pg_notify(
                'docrequestassignee',
                jsonb_build_object(
                    'Request', NEW,
                    'User', action_user
                ) #>> '{}');
        ELSIF OLD.id IS NOT NULL AND COALESCE(NEW.completion_time, '-infinity'::TIMESTAMPTZ) != COALESCE(OLD.completion_time, 'infinity'::TIMESTAMPTZ) THEN
            PERFORM pg_notify(
                'docrequeststatus',
                jsonb_build_object(
                    'Request', NEW,
                    'User', action_user
                ) #>> '{}');
        END IF;

        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS document_request_change_notification ON document_requests;
CREATE TRIGGER document_request_change_notification
    AFTER INSERT OR UPDATE ON document_requests
    FOR EACH ROW
    EXECUTE FUNCTION trigger_document_request_change_notification();
