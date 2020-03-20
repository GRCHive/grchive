DROP TABLE IF EXISTS sql_request_comment_threads;
CREATE TABLE sql_request_comment_threads (
    sql_request_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    thread_id BIGINT NOT NULL,
    FOREIGN KEY(thread_id) REFERENCES comment_threads(id) ON DELETE CASCADE,
    FOREIGN KEY(sql_request_id, org_id) REFERENCES database_sql_query_requests(id, org_id) ON DELETE CASCADE,
    UNIQUE(sql_request_id, org_id, thread_id)
);

CREATE OR REPLACE FUNCTION convert_old_sql_request_thread(input_sql_request_id BIGINT, input_org_id INTEGER)
    RETURNS VOID AS
$$
    DECLARE
        new_thread_id BIGINT;
    BEGIN
        INSERT INTO comment_threads(org_id)
        VALUES (input_org_id)
        RETURNING id INTO new_thread_id;

        INSERT INTO sql_request_comment_threads(sql_request_id, org_id, thread_id)
        VALUES (input_sql_request_id, input_org_id, new_thread_id);

        UPDATE comments
        SET thread_id = new_thread_id
        FROM sql_request_comments AS req
        WHERE req.comment_id = id
            AND req.sql_request_id = input_sql_request_id
            AND req.org_id = input_org_id;
    END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    PERFORM convert_old_sql_request_thread(ut.sql_request_id, ut.org_id)
    FROM (
        SELECT DISTINCT sql_request_id, org_id
        FROM sql_request_comments
    ) AS ut;
END $$ LANGUAGE plpgsql;

DROP FUNCTION IF EXISTS convert_old_sql_request_thread;

CREATE OR REPLACE FUNCTION sql_request_comment_threads_cleanup()
    RETURNS trigger AS
$$
    BEGIN
        DELETE FROM comment_threads
        WHERE id = OLD.thread_id;
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_delete_sql_request_comment_threads
    AFTER DELETE ON sql_request_comment_threads
    FOR EACH ROW
    EXECUTE FUNCTION sql_request_comment_threads_cleanup();
