CREATE OR REPLACE FUNCTION tmp_fn_backfill(input_file_id BIGINT, input_org_id INTEGER)
    RETURNS VOID AS
$$
    DECLARE
        new_thread_id BIGINT;
    BEGIN
        INSERT INTO comment_threads(org_id)
        VALUES (input_org_id)
        RETURNING id INTO new_thread_id;

        INSERT INTO file_comment_threads(file_id, org_id, thread_id)
        VALUES (input_file_id, input_org_id, new_thread_id);
    END;
$$ LANGUAGE plpgsql;

DO $$
DECLARE
BEGIN
    PERFORM tmp_fn_backfill(meta.id, meta.org_id)
    FROM file_metadata AS meta
    LEFT JOIN file_comment_threads AS t
        ON t.file_id = meta.id
    WHERE t.file_id IS NULL;
END $$ LANGUAGE plpgsql;

DROP FUNCTION IF EXISTS tmp_fn_backfill;
