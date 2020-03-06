CREATE OR REPLACE FUNCTION audit_database_sql_queries_change(r database_sql_queries, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        event_id BIGINT;
        db_id BIGINT;
    BEGIN
        SELECT m.db_id INTO db_id
        FROM database_sql_metadata AS m
        WHERE m.id = r.metadata_id;

        SELECT generic_audit_event(r.org_id,
            'database_sql_queries',
            r.id,
            jsonb_build_object(
                'sql_metadata_id', r.metadata_id,
                'db_id', db_id),
            action
        ) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(r));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_database_sql_queries_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_database_sql_queries_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_database_sql_queries_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_database_sql_queries_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_database_sql_queries_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_database_sql_queries_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_database_sql_queries
    AFTER INSERT ON database_sql_queries
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_database_sql_queries_change();

CREATE TRIGGER trigger_update_database_sql_queries
    AFTER UPDATE ON database_sql_queries
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_database_sql_queries_change();

CREATE TRIGGER trigger_delete_database_sql_queries
    BEFORE DELETE ON database_sql_queries
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_database_sql_queries_change();
