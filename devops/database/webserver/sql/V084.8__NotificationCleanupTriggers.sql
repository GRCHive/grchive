CREATE OR REPLACE FUNCTION cleanup_users_notifications()
    RETURNS trigger AS
$$
    BEGIN
        DELETE FROM notifications
        WHERE 
            (subject_type = 'users' AND subject_id = OLD.id) OR
            (object_type = 'users' AND object_id = OLD.id) OR
            (indirect_object_type = 'users' AND indirect_object_id = OLD.id);
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_delete_users_notifs ON users;
CREATE TRIGGER trigger_delete_users_notifs 
    BEFORE DELETE ON users
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_users_notifications();

---

CREATE OR REPLACE FUNCTION cleanup_document_requests_notifications()
    RETURNS trigger AS
$$
    BEGIN
        DELETE FROM notifications
        WHERE 
            (subject_type = 'document_requests' AND subject_id = OLD.id) OR
            (object_type = 'document_requests' AND object_id = OLD.id) OR
            (indirect_object_type = 'document_requests' AND indirect_object_id = OLD.id);
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_delete_document_requests_notifs ON document_requests;
CREATE TRIGGER trigger_delete_document_requests_notifs 
    BEFORE DELETE ON document_requests
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_document_requests_notifications();

---

CREATE OR REPLACE FUNCTION cleanup_database_sql_query_requests_notifications()
    RETURNS trigger AS
$$
    BEGIN
        DELETE FROM notifications
        WHERE 
            (subject_type = 'database_sql_query_requests' AND subject_id = OLD.id) OR
            (object_type = 'database_sql_query_requests' AND object_id = OLD.id) OR
            (indirect_object_type = 'database_sql_query_requests' AND indirect_object_id = OLD.id);
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_delete_database_sql_query_requests_notifs ON database_sql_query_requests;
CREATE TRIGGER trigger_delete_database_sql_query_requests_notifs 
    BEFORE DELETE ON database_sql_query_requests
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_database_sql_query_requests_notifications();

---

CREATE OR REPLACE FUNCTION cleanup_process_flow_controls_notifications()
    RETURNS trigger AS
$$
    BEGIN
        DELETE FROM notifications
        WHERE 
            (subject_type = 'process_flow_controls' AND subject_id = OLD.id) OR
            (object_type = 'process_flow_controls' AND object_id = OLD.id) OR
            (indirect_object_type = 'process_flow_controls' AND indirect_object_id = OLD.id);
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_delete_process_flow_controls_notifs ON process_flow_controls;
CREATE TRIGGER trigger_delete_process_flow_controls_notifs 
    BEFORE DELETE ON process_flow_controls
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_process_flow_controls_notifications();

---

CREATE OR REPLACE FUNCTION cleanup_file_metadata_notifications()
    RETURNS trigger AS
$$
    BEGIN
        DELETE FROM notifications
        WHERE 
            (subject_type = 'file_metadata' AND subject_id = OLD.id) OR
            (object_type = 'file_metadata' AND object_id = OLD.id) OR
            (indirect_object_type = 'file_metadata' AND indirect_object_id = OLD.id);
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_delete_file_metadata_notifs ON file_metadata;
CREATE TRIGGER trigger_delete_file_metadata_notifs 
    BEFORE DELETE ON file_metadata
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_file_metadata_notifications();
