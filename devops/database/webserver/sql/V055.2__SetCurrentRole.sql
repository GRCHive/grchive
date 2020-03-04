CREATE OR REPLACE FUNCTION set_current_role_for_user_id(userid BIGINT)
    RETURNS void AS
$$
    DECLARE
        rol TEXT;
    BEGIN
        SELECT pg.rolname INTO rol
        FROM pg_roles AS pg
        INNER JOIN postgres_oid_to_users AS lnk
            ON lnk.pg_oid = pg.oid
        WHERE lnk.user_id = userid;

        EXECUTE 'SET LOCAL ROLE ' || rol;
    END;
$$ LANGUAGE plpgsql;
