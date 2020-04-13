CREATE OR REPLACE FUNCTION trigger_sql_request_change_notification()
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
                'sqlrequestassignee',
                jsonb_build_object(
                    'Request', NEW,
                    'OldRequest', OLD,
                    'User', action_user
                ) #>> '{}');
        END IF;

        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
