DO $$ BEGIN
    PERFORM fix_users(pg.rolname)
    FROM postgres_oid_to_users AS lnk
    INNER JOIN pg_roles AS pg
        ON pg.oid = lnk.pg_oid;
END $$ LANGUAGE plpgsql;
