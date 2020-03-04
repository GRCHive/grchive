CREATE OR REPLACE FUNCTION audit_database_sql_query_requests_change(r database_sql_query_requests, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        event_id BIGINT;
    BEGIN
        SELECT generic_audit_event(r.org_id,
            'database_sql_query_requests',
            r.id,
            jsonb_build_object('query_id', r.query_id),
            action
        ) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(r));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_database_sql_query_requests_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_database_sql_query_requests_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_database_sql_query_requests_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_database_sql_query_requests_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_database_sql_query_requests_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_database_sql_query_requests_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_database_sql_query_requests
    AFTER INSERT ON database_sql_query_requests
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_database_sql_query_requests_change();

CREATE TRIGGER trigger_update_database_sql_query_requests
    AFTER UPDATE ON database_sql_query_requests
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_database_sql_query_requests_change();

CREATE TRIGGER trigger_delete_database_sql_query_requests
    BEFORE DELETE ON database_sql_query_requests
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_database_sql_query_requests_change();
