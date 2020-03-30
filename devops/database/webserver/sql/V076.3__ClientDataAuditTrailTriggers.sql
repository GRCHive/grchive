CREATE OR REPLACE FUNCTION audit_client_data_change(r client_data, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        cat_id BIGINT;
        event_id BIGINT;
    BEGIN
        SELECT generic_audit_event(r.org_id,
            'client_data',
            r.id,
            '{}'::jsonb,
            action
        ) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(r));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_client_data_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_client_data_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_client_data_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_client_data_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_client_data_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_client_data_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_client_data
    AFTER INSERT ON client_data
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_client_data_change();

CREATE TRIGGER trigger_update_client_data
    AFTER UPDATE ON client_data
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_client_data_change();

CREATE TRIGGER trigger_delete_client_data
    BEFORE DELETE ON client_data
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_client_data_change();
