CREATE OR REPLACE FUNCTION trigger_generic_request_change_notification()
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

        IF NEW.assignee != COALESCE(OLD.assignee, -1) OR NEW.assignee IS NULL THEN
            PERFORM pg_notify(
                'genrequestassignee',
                jsonb_build_object(
                    'Request', NEW,
                    'OldRequest', OLD,
                    'User', action_user
                ) #>> '{}');
        END IF;

        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS generic_request_change_notification ON generic_requests;
CREATE TRIGGER generic_request_change_notification
    AFTER INSERT OR UPDATE ON generic_requests
    FOR EACH ROW
    EXECUTE FUNCTION trigger_generic_request_change_notification();

CREATE OR REPLACE FUNCTION trigger_generic_request_approve_change_notification()
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

        PERFORM pg_notify(
            'genrequeststatus',
            jsonb_build_object(
                'Approval', NEW,
                'User', action_user
            ) #>> '{}');

        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS generic_request_approve_change_notification ON generic_requests;
CREATE TRIGGER generic_request_approve_change_notification
    AFTER INSERT ON generic_requests
    FOR EACH ROW
    EXECUTE FUNCTION trigger_generic_request_approve_change_notification();
