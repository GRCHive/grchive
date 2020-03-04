CREATE OR REPLACE FUNCTION audit_vendors_change(r vendors, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        event_id BIGINT;
    BEGIN
        SELECT generic_audit_event(r.org_id,
            'vendors',
            r.id,
            '{}'::jsonb,
            action
        ) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(r));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_vendors_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_vendors_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_vendors_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_vendors_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_vendors_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_vendors_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_vendors
    AFTER INSERT ON vendors
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_vendors_change();

CREATE TRIGGER trigger_update_vendors
    AFTER UPDATE ON vendors
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_vendors_change();

CREATE TRIGGER trigger_delete_vendors
    BEFORE DELETE ON vendors
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_vendors_change();
