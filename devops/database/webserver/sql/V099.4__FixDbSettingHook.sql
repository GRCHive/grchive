CREATE OR REPLACE FUNCTION sync_database_settings()
    RETURNS trigger AS
$$
    BEGIN
        INSERT INTO database_settings (db_id, org_id)
        VALUES (NEW.id, NEW.org_id);
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
