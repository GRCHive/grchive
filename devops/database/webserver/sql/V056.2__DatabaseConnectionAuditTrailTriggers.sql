CREATE OR REPLACE FUNCTION audit_database_connection_change(conn database_connection_info, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        event_id BIGINT;
    BEGIN
        SELECT generic_audit_event(conn.org_id, 'database_connection_info', conn.id, jsonb_build_object('db_id', conn.db_id), action) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(conn));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_database_connection_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_database_connection_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_database_connection_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_database_connection_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_database_connection_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_database_connection_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_database_connection_info
    AFTER INSERT ON database_connection_info
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_database_connection_change();

CREATE TRIGGER trigger_update_database_connection_info
    AFTER UPDATE ON database_connection_info
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_database_connection_change();

CREATE TRIGGER trigger_delete_database_connection_info
    BEFORE DELETE ON database_connection_info
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_database_connection_change();
