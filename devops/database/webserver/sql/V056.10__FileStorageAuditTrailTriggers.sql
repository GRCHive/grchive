CREATE OR REPLACE FUNCTION audit_file_storage_change(r file_storage, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        cat_id BIGINT;
        event_id BIGINT;
    BEGIN
        SELECT f.category_id INTO cat_id
        FROM file_metadata AS f
        INNER JOIN file_storage AS s
            ON s.metadata_id = f.id
        WHERE s.id = r.id;

        SELECT generic_audit_event(r.org_id,
            'file_storage',
            r.id,
            jsonb_build_object(
                'file_id', r.metadata_id,
                'cat_id', cat_id),
            action
        ) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(r));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_file_storage_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_file_storage_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_file_storage_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_file_storage_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_file_storage_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_file_storage_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_file_storage
    AFTER INSERT ON file_storage
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_file_storage_change();

CREATE TRIGGER trigger_update_file_storage
    AFTER UPDATE ON file_storage
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_file_storage_change();

CREATE TRIGGER trigger_delete_file_storage
    BEFORE DELETE ON file_storage
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_file_storage_change();
