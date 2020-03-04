CREATE OR REPLACE FUNCTION audit_database_change(db database_resources, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        event_id BIGINT;
    BEGIN
        SELECT generic_audit_event(db.org_id, 'database_resources', db.id, '{}'::jsonb, action) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(db));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_database_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_database_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_database_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_database_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_database_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_database_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_database_resources
    AFTER INSERT ON database_resources
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_database_change();

CREATE TRIGGER trigger_update_database_resources
    AFTER UPDATE ON database_resources
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_database_change();

CREATE TRIGGER trigger_delete_database_resources
    BEFORE DELETE ON database_resources
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_database_change();
