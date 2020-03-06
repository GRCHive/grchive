CREATE OR REPLACE FUNCTION add_postgres_role_for_user(u users)
    RETURNS VOID AS
$$
    DECLARE
        parent_role TEXT;
        nm TEXT;
    BEGIN
        SELECT current_user INTO parent_role;
        SELECT create_pg_role_name_for_user(u) INTO nm;

        -- It's probably not super great to create SUPERUSER
        -- roles...but we're already logging in using a SUPERUSER
        -- role and this user can't login anyway.
        EXECUTE 'CREATE ROLE ' || nm || ' WITH INHERIT NOLOGIN PASSWORD NULL';
        EXECUTE 'GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ' || nm;
        EXECUTE 'GRANT ' || nm || ' TO "' || parent_role || '"';

        INSERT INTO postgres_oid_to_users(pg_oid, user_id)
        SELECT rl.oid, u.id
        FROM pg_roles AS rl
        WHERE rl.rolname = nm;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION fix_users(u TEXT)
    RETURNS VOID AS
$$
    DECLARE
        parent_role TEXT;
    BEGIN
        SELECT current_user INTO parent_role;
        EXECUTE 'GRANT ' || u || ' TO "' || parent_role || '"';
    END;
$$ LANGUAGE plpgsql;

DO $$ BEGIN
    PERFORM fix_users(pg.rolname)
    FROM postgres_oid_to_users AS lnk
    INNER JOIN pg_roles AS pg
        ON pg.oid = lnk.pg_oid;
END $$ LANGUAGE plpgsql;
