CREATE OR REPLACE FUNCTION sync_database_settings()
    RETURNS trigger AS
$$
    BEGIN
        INSERT INTO database_settings (db_id, org_id)
        VALUES (NEW.id, NEW.org_id);
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_sync_database_settings ON database_resources;
DROP TRIGGER IF EXISTS trigger_sync_database_settings ON users;
CREATE TRIGGER trigger_sync_database_settings
    AFTER INSERT ON database_resources
    FOR EACH ROW
    EXECUTE FUNCTION sync_database_settings();

INSERT INTO database_settings (db_id, org_id)
SELECT id, org_id
FROM database_resources
ON CONFLICT (db_id, org_id) DO NOTHING;
