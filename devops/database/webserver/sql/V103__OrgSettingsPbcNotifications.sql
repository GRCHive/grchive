CREATE TABLE org_pbc_notification_cadence_settings (
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    days_before_due INTEGER NOT NULL
);

CREATE INDEX ON org_pbc_notification_cadence_settings(org_id);

-- Default notification cadence is day of, 1 day before, 3 days before, 7 days before, 14 days
CREATE OR REPLACE FUNCTION create_org_default_pbc_notification_cadence(input_org_id INTEGER)
    RETURNS VOID AS
$$
    BEGIN
        INSERT INTO org_pbc_notification_cadence_settings (org_id, days_before_due)
        VALUES
            (input_org_id, 0),
            (input_org_id, 1),
            (input_org_id, 3),
            (input_org_id, 7),
            (input_org_id, 14);
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION fn_tg_org_pbc_notification_cadence()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM create_org_default_pbc_notification_cadence(NEW.id);
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_org_pbc_notification_cadence
    AFTER INSERT ON organizations
    FOR EACH ROW
    EXECUTE FUNCTION fn_tg_org_pbc_notification_cadence();

DO $$ BEGIN
    PERFORM fix_users(pg.rolname)
    FROM postgres_oid_to_users AS lnk
    INNER JOIN pg_roles AS pg
        ON pg.oid = lnk.pg_oid;

    PERFORM create_org_default_pbc_notification_cadence(id)
    FROM organizations;
END $$ LANGUAGE plpgsql;
