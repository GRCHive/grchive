CREATE OR REPLACE FUNCTION trigger_control_change_notification()
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

        IF NEW.owner_id != COALESCE(OLD.owner_id, -1) THEN
            PERFORM pg_notify(
                'controlowner',
                jsonb_build_object(
                    'Control', NEW,
                    'User', action_user
                ) #>> '{}');
        END IF;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS control_change_notification ON process_flow_controls;
CREATE TRIGGER control_change_notification
    AFTER INSERT OR UPDATE ON process_flow_controls
    FOR EACH ROW
    EXECUTE FUNCTION trigger_control_change_notification();
