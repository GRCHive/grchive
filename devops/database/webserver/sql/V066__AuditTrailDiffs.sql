CREATE OR REPLACE FUNCTION jsonb_diff(a jsonb, b jsonb)
    RETURNS jsonb AS
$$
    DECLARE
        ret jsonb := '{}'::jsonb;
        merged_keys text[];
        key text;
        a_val text;
        b_val text;
    BEGIN
        WITH union_query AS (
            SELECT jsonb_object_keys(a)
            UNION SELECT jsonb_object_keys(b)
        )
        SELECT ARRAY_AGG(jsonb_object_keys)
        FROM union_query
        INTO merged_keys;

        FOREACH key IN ARRAY merged_keys
        LOOP
            SELECT a ->> key INTO a_val;
            SELECT b ->> key INTO b_val;

            IF a_val != b_val THEN
                ret = ret || jsonb_build_object(key, 
                    jsonb_build_object(
                        'new', b ->> key,
                        'old', a ->> key
                    )
                );
            END IF;
        END LOOP;

        RETURN ret;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION create_modification_diff(event_id BIGINT)
    RETURNS jsonb AS
$$
    DECLARE
        current_event global_audit_event_history;
        previous_event global_audit_event_history;
    BEGIN
        SELECT *
        FROM global_audit_event_history
        WHERE id = event_id
        INTO current_event;

        SELECT *
        FROM global_audit_event_history
        WHERE id < event_id
            AND resource_type = current_event.resource_type
            AND resource_id = current_event.resource_id
            AND org_id = current_event.org_id
            AND action != 'SELECT'
        INTO previous_event
        ORDER BY id DESC
        LIMIT 1;

        IF previous_event IS NOT NULL THEN
            RETURN jsonb_diff(
                (SELECT data FROM audit_resource_modifications AS m WHERE m.event_id = previous_event.id),
                (SELECT data FROM audit_resource_modifications AS m WHERE m.event_id = current_event.id)
            );
        END IF;

        RETURN '{}'::jsonb;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE VIEW audit_resource_modification_diffs AS
    SELECT
        m.event_id AS event_id,
        create_modification_diff(m.event_id) AS diff
    FROM audit_resource_modifications AS m
    INNER JOIN global_audit_event_history AS hist
        ON hist.id = m.event_id
    WHERE hist.action = 'UPDATE';
