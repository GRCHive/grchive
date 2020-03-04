CREATE OR REPLACE FUNCTION audit_infrastructure_change(sv infrastructure_servers, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        event_id BIGINT;
    BEGIN
        SELECT generic_audit_event(sv.org_id, 'infrastructure_servers', sv.id, '{}'::jsonb, action) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(sv));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_infrastructure_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_infrastructure_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_infrastructure_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_infrastructure_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_infrastructure_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_infrastructure_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_infrastructure_servers
    AFTER INSERT ON infrastructure_servers
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_infrastructure_change();

CREATE TRIGGER trigger_update_infrastructure_servers
    AFTER UPDATE ON infrastructure_servers
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_infrastructure_change();

CREATE TRIGGER trigger_delete_infrastructure_servers
    BEFORE DELETE ON infrastructure_servers
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_infrastructure_change();
