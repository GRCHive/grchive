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
        FROM code_to_client_scripts AS lnk
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

CREATE OR REPLACE FUNCTION insert_audit_managed_code_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_managed_code_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_managed_code_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_managed_code_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_managed_code_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_managed_code_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_managed_code
    AFTER INSERT ON managed_code
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_managed_code_change();

CREATE TRIGGER trigger_update_managed_code
    AFTER UPDATE ON managed_code
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_managed_code_change();

CREATE TRIGGER trigger_delete_managed_code
    BEFORE DELETE ON managed_code
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_managed_code_change();
