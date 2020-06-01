CREATE OR REPLACE FUNCTION fn_shell_script_delete_on_generic_request_delete()
    RETURNS trigger AS
$$
    BEGIN
        DELETE FROM shell_script_runs AS ssr
        USING request_to_shell_run_link AS srl,
              generic_approval AS ga
        WHERE srl.run_id = ssr.id AND srl.request_id = OLD.id AND NOT ga.response;
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS shell_script_delete_on_generic_request_delete ON generic_requests;
CREATE TRIGGER shell_script_delete_on_generic_request_delete
    BEFORE DELETE ON generic_requests
    FOR EACH ROW
    EXECUTE FUNCTION fn_shell_script_delete_on_generic_request_delete();

DO $$ BEGIN
    PERFORM fix_users(pg.rolname)
    FROM postgres_oid_to_users AS lnk
    INNER JOIN pg_roles AS pg
        ON pg.oid = lnk.pg_oid;
END $$ LANGUAGE plpgsql;
