CREATE OR REPLACE FUNCTION regrant_access_to_users(u TEXT)
    RETURNS VOID AS
$$
    DECLARE
        parent_role TEXT;
    BEGIN
        SELECT current_user INTO parent_role;
        EXECUTE 'GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ' || u;
        EXECUTE 'GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ' || u;
        EXECUTE 'GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO ' || u;
        EXECUTE 'GRANT ALL PRIVILEGES ON ALL PROCEDURES IN SCHEMA public TO ' || u;
        EXECUTE 'GRANT ALL PRIVILEGES ON ALL ROUTINES IN SCHEMA public TO ' || u;
    END;
$$ LANGUAGE plpgsql;

DO $$ BEGIN
    PERFORM regrant_access_to_users(pg.rolname)
    FROM postgres_oid_to_users AS lnk
    INNER JOIN pg_roles AS pg
        ON pg.oid = lnk.pg_oid;
END $$ LANGUAGE plpgsql;
