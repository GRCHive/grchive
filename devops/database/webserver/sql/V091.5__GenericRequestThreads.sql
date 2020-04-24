CREATE TABLE generic_request_comment_threads (
    request_id BIGINT NOT NULL UNIQUE REFERENCES generic_requests(id) ON DELETE CASCADE,
    thread_id BIGINT NOT NULL UNIQUE REFERENCES comment_threads(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION sync_generic_request_comment_thread()
    RETURNS trigger AS
$$
    DECLARE
        new_thread_id BIGINT;
    BEGIN
        INSERT INTO comment_threads(org_id)
        VALUES (NEW.org_id)
        RETURNING id INTO new_thread_id;

        INSERT INTO generic_request_comment_threads(request_id, thread_id)
        VALUES (NEW.id, new_thread_id);
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_sync_generic_request_comment_thread ON generic_request;
CREATE TRIGGER trigger_sync_generic_request_comment_thread
    AFTER INSERT ON generic_requests
    FOR EACH ROW
    EXECUTE FUNCTION sync_generic_request_comment_thread();

CREATE OR REPLACE FUNCTION tmp_fn_backfill(input_request_id BIGINT, input_org_id INTEGER)
    RETURNS VOID AS
$$
    DECLARE
        new_thread_id BIGINT;
    BEGIN
        INSERT INTO comment_threads(org_id)
        VALUES (input_org_id)
        RETURNING id INTO new_thread_id;

        INSERT INTO generic_request_comment_threads(request_id, thread_id)
        VALUES (input_request_id, new_thread_id);
    END;
$$ LANGUAGE plpgsql;

DO $$
DECLARE
BEGIN
    PERFORM tmp_fn_backfill(req.id, req.org_id)
    FROM generic_requests AS req;
END $$ LANGUAGE plpgsql;

DROP FUNCTION IF EXISTS tmp_fn_backfill;
