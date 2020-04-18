CREATE OR REPLACE FUNCTION audit_managed_code_change(r managed_code, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        event_id BIGINT;
        linked_data BIGINT;
        linked_script BIGINT;
    BEGIN
        SELECT generic_audit_event(r.org_id,
            'managed_code',
            r.id,
            '{}'::jsonb,
            action
        ) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(r));
    END;
$$ LANGUAGE plpgsql;
