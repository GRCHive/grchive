DROP TRIGGER IF EXISTS trigger_add_postgres_role_for_user ON users;
DROP TRIGGER IF EXISTS trigger_delete_postgres_role_for_user ON users;
DROP TABLE IF EXISTS postgres_oid_to_users;

CREATE TABLE postgres_oid_to_users (
    pg_oid oid NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION create_pg_role_name_for_user(u users)
    RETURNS TEXT AS
$$
    BEGIN
        RETURN 'user_' || u.id;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_postgres_role_for_user(u users)
    RETURNS VOID AS
$$
    DECLARE
        nm TEXT;
    BEGIN
        SELECT create_pg_role_name_for_user(u) INTO nm;

        -- It's probably not super great to create SUPERUSER
        -- roles...but we're already logging in using a SUPERUSER
        -- role and this user can't login anyway.
        EXECUTE 'CREATE ROLE ' || nm || ' WITH INHERIT NOLOGIN PASSWORD NULL';
        EXECUTE 'GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ' || nm;

        INSERT INTO postgres_oid_to_users(pg_oid, user_id)
        SELECT rl.oid, u.id
        FROM pg_roles AS rl
        WHERE rl.rolname = nm;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION on_user_create()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM add_postgres_role_for_user(NEW);
        return NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_add_postgres_role_for_user
    AFTER INSERT ON users
    FOR EACH ROW
    EXECUTE FUNCTION on_user_create();

CREATE OR REPLACE FUNCTION delete_postgres_role_for_user(u users)
    RETURNS VOID AS
$$
    DECLARE
        nm TEXT;
    BEGIN
        SELECT create_pg_role_name_for_user(u) INTO nm;
        DELETE FROM postgres_oid_to_users
        WHERE user_id = u.id;
        DROP ROLE nm;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION on_user_delete()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM delete_postgres_role_for_user(OLD);
        return OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_delete_postgres_role_for_user
    BEFORE DELETE ON users
    FOR EACH ROW
    EXECUTE FUNCTION on_user_delete();

DO $$ BEGIN
    PERFORM add_postgres_role_for_user(u.*)
    FROM users AS u;
END $$ LANGUAGE plpgsql;
