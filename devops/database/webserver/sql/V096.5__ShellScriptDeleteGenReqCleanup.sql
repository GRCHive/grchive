CREATE OR REPLACE FUNCTION fn_shell_script_delete_on_generic_request_delete()
    RETURNS trigger AS
$$
    BEGIN
        DELETE FROM shell_script_runs AS ssr
        USING request_to_shell_run_link AS srl
        WHERE srl.run_id = ssr.id AND request_id = OLD.id;
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS shell_script_delete_on_generic_request_delete ON generic_requests;
CREATE TRIGGER shell_script_delete_on_generic_request_delete
    AFTER DELETE ON generic_requests
    FOR EACH ROW
    EXECUTE FUNCTION fn_shell_script_delete_on_generic_request_delete();
