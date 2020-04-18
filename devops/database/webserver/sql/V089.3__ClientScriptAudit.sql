CREATE OR REPLACE FUNCTION audit_client_scripts_change(r client_scripts, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        event_id BIGINT;
    BEGIN
        SELECT generic_audit_event(r.org_id,
            'client_scripts',
            r.id,
            '{}'::jsonb,
            action
        ) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(r));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_client_scripts_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_client_scripts_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_client_scripts_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_client_scripts_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_client_scripts_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_client_scripts_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_client_scripts
    AFTER INSERT ON client_scripts
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_client_scripts_change();

CREATE TRIGGER trigger_update_client_scripts
    AFTER UPDATE ON client_scripts
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_client_scripts_change();

CREATE TRIGGER trigger_delete_client_scripts
    BEFORE DELETE ON client_scripts
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_client_scripts_change();
