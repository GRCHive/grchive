CREATE OR REPLACE FUNCTION generic_audit_event(
    org_id INTEGER,
    resource_type VARCHAR(64),
    resource_id BIGINT,
    resource_extra_data JSONB,
    action VARCHAR(64)
)
    RETURNS BIGINT AS
$$
    DECLARE
        event_id BIGINT;
        pgrole_id oid;
    BEGIN
        SELECT oid INTO pgrole_id
        FROM pg_roles
        WHERE rolname = current_user;

        INSERT INTO global_audit_event_history(
            org_id,
            resource_type,
            resource_id,
            resource_extra_data,
            action,
            performed_at,
            pgrole_id
        )
        VALUES (
            org_id,
            resource_type,
            CAST(resource_id AS VARCHAR(128)),
            resource_extra_data,
            action,
            NOW(),
            pgrole_id
        )
        RETURNING id INTO event_id;
        RETURN event_id;
    END;
$$ LANGUAGE plpgsql;
