CREATE TABLE database_settings (
    db_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    auto_refresh_task BIGINT REFERENCES scheduled_tasks(id) ON DELETE SET NULL,
    FOREIGN KEY(db_id, org_id) REFERENCES database_resources(id, org_id) ON DELETE CASCADE,
    UNIQUE(db_id, org_id)
);

INSERT INTO database_settings (db_id, org_id)
SELECT id, org_id
FROM database_resources;
