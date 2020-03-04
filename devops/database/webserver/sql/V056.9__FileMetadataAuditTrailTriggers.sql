CREATE OR REPLACE FUNCTION audit_file_metadata_change(r file_metadata, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        event_id BIGINT;
    BEGIN
        SELECT generic_audit_event(r.org_id,
            'file_metadata',
            r.id,
            jsonb_build_object('cat_id', r.category_id),
            action
        ) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(r));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_file_metadata_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_file_metadata_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_file_metadata_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_file_metadata_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_file_metadata_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_file_metadata_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_file_metadata
    AFTER INSERT ON file_metadata
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_file_metadata_change();

CREATE TRIGGER trigger_update_file_metadata
    AFTER UPDATE ON file_metadata
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_file_metadata_change();

CREATE TRIGGER trigger_delete_file_metadata
    BEFORE DELETE ON file_metadata
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_file_metadata_change();
