CREATE OR REPLACE FUNCTION audit_managed_code_change(r managed_code, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        event_id BIGINT;
        linked_data BIGINT;
        linked_script BIGINT;
    BEGIN
        SELECT lnk.data_id INTO linked_data
        FROM code_to_client_data_link AS lnk
        WHERE lnk.code_id = r.id;

        SELECT lnk.script_id INTO linked_script
        FROM code_to_client_scripts_link AS lnk
        WHERE lnk.code_id = r.id;

        SELECT generic_audit_event(r.org_id,
            'managed_code',
            r.id,
            jsonb_build_object(
                'client_data_id', linked_data,
                'client_script_id', linked_script
            ),
            action
        ) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(r));
    END;
$$ LANGUAGE plpgsql;
