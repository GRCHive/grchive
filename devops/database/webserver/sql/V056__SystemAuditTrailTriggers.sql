CREATE OR REPLACE FUNCTION audit_system_change(sys systems, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        event_id BIGINT;
    BEGIN
        SELECT generic_audit_event(sys.org_id, 'systems', sys.id, '{}'::jsonb, action) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(sys));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_system_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_system_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_system_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_system_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_system_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_system_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_systems
    AFTER INSERT ON systems
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_system_change();

CREATE TRIGGER trigger_update_systems
    AFTER UPDATE ON systems
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_system_change();

CREATE TRIGGER trigger_delete_systems
    BEFORE DELETE ON systems
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_system_change();
