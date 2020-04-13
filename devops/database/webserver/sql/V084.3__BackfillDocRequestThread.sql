CREATE OR REPLACE FUNCTION tmp_fn_backfill(input_request_id BIGINT, input_org_id INTEGER)
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
    END;
$$ LANGUAGE plpgsql;

DO $$
DECLARE
BEGIN
    PERFORM tmp_fn_backfill(doc.id, doc.org_id)
    FROM document_requests AS doc
    LEFT JOIN document_request_comment_threads AS t
        ON t.request_id = doc.id
    WHERE t.request_id IS NULL;
END $$ LANGUAGE plpgsql;

DROP FUNCTION IF EXISTS tmp_fn_backfill;
