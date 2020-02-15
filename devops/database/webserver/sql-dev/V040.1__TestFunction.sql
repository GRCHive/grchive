DROP FUNCTION IF EXISTS test_grab_all_user_org_controls;
CREATE OR REPLACE FUNCTION test_grab_all_user_org_controls(userId BIGINT, orgId INTEGER)
    RETURNS SETOF process_flow_controls AS
$$
    BEGIN
        return QUERY SELECT *
        FROM process_flow_controls
        WHERE owner_id = userId AND org_id = orgId;
    END;
$$ LANGUAGE plpgsql;
