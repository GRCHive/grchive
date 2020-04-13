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

        IF (NEW.owner_id != COALESCE(OLD.owner_id, -1)) OR NEW.owner_id IS NULL THEN
            PERFORM pg_notify(
                'controlowner',
                jsonb_build_object(
                    'Control', NEW,
                    'OldControl', OLD,
                    'User', action_user
                ) #>> '{}');
        END IF;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
