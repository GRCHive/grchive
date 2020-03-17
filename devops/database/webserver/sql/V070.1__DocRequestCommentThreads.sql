DROP TABLE IF EXISTS document_request_comment_threads;
CREATE TABLE document_request_comment_threads (
    request_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    thread_id BIGINT NOT NULL,
    FOREIGN KEY(thread_id) REFERENCES comment_threads(id) ON DELETE CASCADE,
    FOREIGN KEY(request_id, org_id) REFERENCES document_requests(id, org_id) ON DELETE CASCADE,
    UNIQUE(request_id, org_id, thread_id)
);

CREATE OR REPLACE FUNCTION convert_old_doc_request_thread(input_request_id BIGINT, input_org_id INTEGER)
    RETURNS VOID AS
$$
    DECLARE
        new_thread_id BIGINT;
    BEGIN
        INSERT INTO comment_threads(org_id)
        VALUES (input_org_id)
        RETURNING id INTO new_thread_id;

        INSERT INTO document_request_comment_threads(request_id, org_id, thread_id)
        VALUES (input_request_id, input_org_id, new_thread_id);

        UPDATE comments
        SET thread_id = new_thread_id
        FROM document_request_comments AS req
        WHERE req.comment_id = id
            AND req.request_id = input_request_id
            AND req.org_id = input_org_id;
    END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    PERFORM (WITH unique_threads AS (
        SELECT DISTINCT request_id, org_id
        FROM document_request_comments
    )
    SELECT convert_old_doc_request_thread(request_id, org_id)
    FROM unique_threads);
END $$ LANGUAGE plpgsql;

DROP FUNCTION IF EXISTS convert_old_doc_request_thread;

CREATE OR REPLACE FUNCTION document_request_comment_threads_cleanup()
    RETURNS trigger AS
$$
    BEGIN
        DELETE FROM comment_threads
        WHERE id = OLD.thread_id;
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_delete_document_request_comment_threads
    AFTER DELETE ON document_request_comment_threads
    FOR EACH ROW
    EXECUTE FUNCTION document_request_comment_threads_cleanup();
