CREATE OR REPLACE FUNCTION sync_document_requests_comment_thread()
    RETURNS trigger AS
$$
    DECLARE
        new_thread_id BIGINT;
    BEGIN
        INSERT INTO comment_threads(org_id)
        VALUES (NEW.org_id)
        RETURNING id INTO new_thread_id;

        INSERT INTO document_request_comment_threads(request_id, org_id, thread_id)
        VALUES (NEW.id, NEW.org_id, new_thread_id);
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_sync_document_requests_comment_thread ON document_requests;
CREATE TRIGGER trigger_sync_document_requests_comment_thread
    AFTER INSERT ON document_requests
    FOR EACH ROW
    EXECUTE FUNCTION sync_document_requests_comment_thread();

CREATE OR REPLACE FUNCTION sync_file_metadata_comment_thread()
    RETURNS trigger AS
$$
    DECLARE
        new_thread_id BIGINT;
    BEGIN
        INSERT INTO comment_threads(org_id)
        VALUES (NEW.org_id)
        RETURNING id INTO new_thread_id;

        INSERT INTO file_comment_threads(file_id, org_id, thread_id)
        VALUES (NEW.id, NEW.org_id, new_thread_id);
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_sync_file_metadata_comment_thread ON file_metadata;
CREATE TRIGGER trigger_sync_file_metadata_comment_thread
    AFTER INSERT ON file_metadata
    FOR EACH ROW
    EXECUTE FUNCTION sync_file_metadata_comment_thread();

CREATE OR REPLACE FUNCTION sync_sql_request_comment_thread()
    RETURNS trigger AS
$$
    DECLARE
        new_thread_id BIGINT;
    BEGIN
        INSERT INTO comment_threads(org_id)
        VALUES (NEW.org_id)
        RETURNING id INTO new_thread_id;

        INSERT INTO sql_request_comment_threads(sql_request_id, org_id, thread_id)
        VALUES (NEW.id, NEW.org_id, new_thread_id);
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_sync_sql_request_comment_thread ON database_sql_query_requests;
CREATE TRIGGER trigger_sync_sql_request_comment_thread
    AFTER INSERT ON database_sql_query_requests
    FOR EACH ROW
    EXECUTE FUNCTION sync_sql_request_comment_thread();
