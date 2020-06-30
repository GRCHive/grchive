CREATE OR REPLACE FUNCTION create_org_default_pbc_notification_cadence(input_org_id INTEGER)
    RETURNS VOID AS
$$
    BEGIN
        INSERT INTO org_pbc_notification_cadence_settings (org_id, days_before_due)
        VALUES
            (input_org_id, -7),
            (input_org_id, -3),
            (input_org_id, -1),
            (input_org_id, 0),
            (input_org_id, 1),
            (input_org_id, 3),
            (input_org_id, 7),
            (input_org_id, 14);
    END;
$$ LANGUAGE plpgsql;
